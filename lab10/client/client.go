package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
)

var done chan interface{}
var interrupt chan os.Signal

func handleReceive(conn *websocket.Conn) {
	defer close(done)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Print(msg)
	}
}

func main() {
	done = make(chan interface{})
	interrupt = make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8282/ftp", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	go handleReceive(c)
	for {
		select {
		case <-done:
			log.Println("connection was closed")
			return
		case <-interrupt:
			err = c.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}
			return
		}
	}
}
