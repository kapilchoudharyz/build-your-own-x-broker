package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/kapilchoudharyz/build-your-own-x/internal/broker"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		b := broker.NewBroker()
		http.HandleFunc("/ws", broker.NewWebSocketHandler(b))
		fmt.Println("Starting an http server...")
		http.ListenAndServe(":8080", nil)
		wg.Done()
	}()
	wg.Wait()
}
