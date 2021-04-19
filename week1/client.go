package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var done chan interface{}
var interrupt chan os.Signal

func receiveHandler(connection *websocket.Conn) {
	defer close(done)

	for {
		_, msg, err := connection.ReadMessage()

		if err != nil {
			log.Println("Error in receive: ", err)
			return
		}

		log.Printf("Received: %s\n", msg)
	}
}

func main() {
	done = make(chan interface{})
	interrupt = make(chan os.Signal)

	signal.Notify(interrupt, os.Interrupt)

	socketUrl := "ws://localhost:8080" + "/predict"

	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)

	if err != nil {
		log.Fatal("Error connecting to server: ", err)
	}

	defer conn.Close()
	go receiveHandler(conn)

	for {
		select {
		case <-time.After(time.Duration(1) * time.Millisecond * 1000):
			err := conn.WriteMessage(websocket.TextMessage, []byte("Hello from the client!"))

			if err != nil {
				log.Println("Error during writing to socket: ", err)
				return
			}
		case <-interrupt:
			log.Println("Received SIGINT interrupt signal, CLosing all pending connections")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing of the socket: ", err)
				return
			}

			select {
			case <-done:
				log.Println("Receiver channel closed! Exiting....")
			case <-time.After(time.Duration(1) * time.Second):
				log.Println("Timeout in closing channel. Exiting....")
			}
			return
		}
	}
}
