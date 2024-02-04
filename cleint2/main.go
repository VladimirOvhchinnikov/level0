package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/stan.go"
)

func main() {
	sc, err := stan.Connect("wbCluster", "client-subscriber", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	// Подписка на канал
	sub, err := sc.Subscribe("wbCluster", func(m *stan.Msg) {
		log.Printf("Received a message: %s\n", string(m.Data))
	}, stan.DeliverAllAvailable())
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	// Ожидание сигнала для завершения
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}
