package random

import (
	"context"
	"fmt"
	"log"

	"github.com/coder/websocket/wsjson"
	"github.com/kapilchoudharyz/build-your-own-x/internal/broker"
)

func NewConnOutHandler() chan broker.WSConnWithID {
	c := make(chan broker.WSConnWithID, 1000)

	return c
}

func Worker(c broker.WSConnWithID) {
	fmt.Print("Inside worker")
	for {
		fmt.Print("Inside worker for loop...")
		var incomingMessage broker.Message[any]
		err := wsjson.Read(context.TODO(), c.Conn, &incomingMessage)
		if err != nil {
			log.Printf("Error reading: %v", err)
			return
		}
		fmt.Printf("REceived %+v\n", incomingMessage)
	}
}

func CreateWorkers(inputChan <-chan broker.WSConnWithID) {
	for c := range inputChan {
		fmt.Printf("wohoo %+v\n", c)
		go func() {
			Worker(c)
		}()
	}
}
