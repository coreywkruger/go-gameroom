package main

import (
	"github.com/gorilla/websocket"
	"log"
)

// Connection -
type Connection struct {
	ID      int
	send    chan []byte
	receive *chan []byte
	closed  *chan int
}

// Reader - reads message from websocket; puts in `receive` channel
func (c *Connection) Reader(ws *websocket.Conn) {
	for {
		_, Message, err := ws.ReadMessage()
		if err != nil {
			log.Println("Couldn't Read:", err.Error())
			break
		}
		*c.receive <- Message
	}
	ws.Close()
}

// Writer - writes message to websocket
func (c *Connection) Writer(ws *websocket.Conn) {
	for message := range c.send {
		err := ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Couldn't Write:", err.Error())
			break
		}
	}
	ws.Close()
}

// Kill - closes channels
func (c *Connection) Kill() {
	close(c.send)
}
