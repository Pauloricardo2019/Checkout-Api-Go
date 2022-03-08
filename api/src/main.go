package main

import (
	"fmt"
	"os"
	"os/signal"
	"ravxcheckout/src/adapter/monitoring"
	"ravxcheckout/src/adapter/nats"
	"ravxcheckout/src/adapter/rest"
	"syscall"
	"time"
)

var chanSignal chan os.Signal

func main() {
	monitoring.StartMonitoring()
	nats.StartSubscriptions()
	rest.InitRest()
	handleExit()
	defer handlePanic()


	/*
		Deixa o sistema em 'loop'.
		Caso essa linha seja apagada, o sistema vai iniciar e finalizar em milésimos de segundo,
	*/
	for {
		time.Sleep(2 * time.Second)
	}
}

func handlePanic() {
	if err := recover(); err != nil {
		chanSignal <- syscall.SIGTERM
		for {
			time.Sleep(2 * time.Second)
		}
	}
}

func handleExit() {
	chanSignal = make(chan os.Signal, 1)
	signal.Notify(chanSignal, os.Interrupt)

	go func() {
		for exitSignal := range chanSignal {
			fmt.Printf("Capturing signal... %v\n", exitSignal)

			fmt.Println("Draining NATS Connection")
			nats.StopSubscriptions()
			monitoring.Flush()

			fmt.Println("Exiting system")
			os.Exit(0)
		}
	}()
}
