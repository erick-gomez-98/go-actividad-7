package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Proceso_ struct {
	Id uint64
	I  uint64
}

type Proceso struct {
	Id uint64
	I  chan uint64
}

func (p *Proceso) start(s uint64) {
	i := uint64(s)
	for {
		time.Sleep(time.Millisecond * 500)
		p.I <- i
		i++
	}
}

func client() {
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	var proceso_ Proceso_
	err = gob.NewDecoder(c).Decode(&proceso_)
	if err != nil {
		fmt.Println(err)
	}

	procesos := []Proceso{}
	p1 := Proceso{Id: proceso_.Id, I: make(chan uint64)}
	procesos = append(procesos, p1)
	go p1.start(proceso_.I)

	cl := make(chan os.Signal)
	signal.Notify(cl, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-cl
		fmt.Println("\r- Proceso terminado")
		p := Proceso_{Id: procesos[0].Id, I: <-procesos[0].I}
		err = gob.NewEncoder(c).Encode(p)
		c.Close()
		os.Exit(0)
	}()

	go func() {
		for {
			time.Sleep(time.Millisecond * 500)
			for _, proceso := range procesos {
				msg := <-proceso.I
				fmt.Printf("%d : %d \n", proceso.Id, msg)
			}
		}
	}()
}

func main() {
	go client()
	var input string
	fmt.Scanln(&input)
}
