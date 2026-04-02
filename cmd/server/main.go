package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/kapilchoudharyz/build-your-own-x/internal/broker"
	"github.com/kapilchoudharyz/build-your-own-x/internal/random"
)

func main() {
	connectionOutputChan := random.NewConnOutHandler()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		b := broker.NewBroker()
		http.HandleFunc("/ws", broker.NewWebSocketHandler(b, connectionOutputChan))
		fmt.Println("Starting an http server...")
		http.ListenAndServe(":8080", nil)
		wg.Done()
	}()
	random.CreateWorkers(connectionOutputChan)
	wg.Wait()
}
