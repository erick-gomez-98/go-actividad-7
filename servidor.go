package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

type Proceso struct {
	Id uint64
	I  chan uint64
}

type Proceso_ struct {
	Id uint64
	I  uint64
}

func (p *Proceso) start(s uint64) {
	i := uint64(s)
	for {
		p.I <- i
		i++
	}
}

func servidor() {
	s, err := net.Listen("tcp", ":9999")
	procesos := []Proceso{}
	p1 := Proceso{Id: 0, I: make(chan uint64)}
	procesos = append(procesos, p1)
	go p1.start(0)
	p2 := Proceso{Id: 1, I: make(chan uint64)}
	procesos = append(procesos, p2)
	go p2.start(0)
	p3 := Proceso{Id: 2, I: make(chan uint64)}
	procesos = append(procesos, p3)
	go p3.start(0)
	p4 := Proceso{Id: 3, I: make(chan uint64)}
	procesos = append(procesos, p4)
	go p4.start(0)
	p5 := Proceso{Id: 4, I: make(chan uint64)}
	procesos = append(procesos, p5)
	go p5.start(0)

	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			time.Sleep(time.Millisecond * 500)
			for _, proceso := range procesos {
				msg := <-proceso.I
				fmt.Printf("%d : %d \n", proceso.Id, msg)
			}
			fmt.Println("-----------")
		}
	}()

	for {
		c, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleClient(c, &procesos)
	}
}

func handleClient(c net.Conn, procesos *[]Proceso) {
	var proceso Proceso
	if len(*procesos) > 0 {
		proceso = (*procesos)[0]
		i := <-proceso.I
		p := Proceso_{Id: proceso.Id, I: i}
		err := gob.NewEncoder(c).Encode(p)
		if err != nil {
			fmt.Println(err)
			return
		}
		remove(procesos, 0)
	}

	var pr Proceso_
	err := gob.NewDecoder(c).Decode(&pr)
	if err == nil {
		p := Proceso{Id: pr.Id, I: make(chan uint64)}
		go p.start(pr.I)
		add(procesos, p)
		return
	}
}

func remove(slice *[]Proceso, s int) {
	ss := *slice
	ss = append(ss[:s], ss[s+1:]...)
	*slice = ss
}

func add(slice *[]Proceso, s Proceso) {
	ss := *slice
	ss = append(ss, s)
	*slice = ss
}

func main() {
	go servidor()

	var input string
	fmt.Scanln(&input)
}



/*
NOTA:
Por los time.sleep que tienen que estar cada 500ms, no pude encontrar una forma 
para que el cliente le regresará el número en el cual va el servidor, funciona todo,
pero solo que este número a la hora de regresarlo viene desfasado por 1 o 2 seg.
*/