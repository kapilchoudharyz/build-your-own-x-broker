package broker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
	"github.com/kapilchoudharyz/build-your-own-x/internal/client"
)

type Broker struct {
	topicClientMap map[string]map[*client.Client]struct{}
}

type Message[T any] struct {
	Topic   string `json:"topic"`
	Payload T      `json:"payload"`
}

type WSConnWithID struct {
	ID   string
	Conn *websocket.Conn
}

func (b *Broker) Subscribe() {
	fmt.Println("Subscribe")
}
func (b *Broker) UnSubscribe() {
	fmt.Println("UnSubscribe")
}
func (b *Broker) Publish() {
	fmt.Println("Publish")
}

func NewBroker() *Broker {
	return &Broker{
		topicClientMap: make(map[string]map[*client.Client]struct{}),
	}
}

func readAndValidateMessages(ctx context.Context, c *websocket.Conn) error {
	for {
		var incomingMessage Message[any]
		err := wsjson.Read(ctx, c, &incomingMessage)
		if err != nil {
			log.Printf("Error reading: %v", err)
			return err
		}
		fmt.Printf("REceived %+v\n", incomingMessage)
	}
}

func NewWebSocketHandler(b *Broker, connectionOutChan chan<- WSConnWithID) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			panic(fmt.Sprintf("Error creating web socket handler: %v\n", err))
		}
		connectionId := uuid.New().String()
		fmt.Println(connectionId)
		defer c.CloseNow()
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
		defer cancel()

		connectionOutChan <- WSConnWithID{
			ID:   connectionId,
			Conn: c,
		}

		// go func() {
		// 	err := readAndValidateMessages(ctx, c)
		// 	if err != nil {
		// 		cancel()
		// 	}
		// }()

		t := time.NewTicker(10 * time.Second)
		for {

			select {
			case <-t.C:
				err = wsjson.Write(ctx, c, "Hello Bhai")
				fmt.Println("Writing over the socket...")
				if err != nil {
					log.Println("err")
					return

				}
			case <-ctx.Done():
				c.Close(websocket.StatusNormalClosure, "")
				return
			}

		}
	}
}
