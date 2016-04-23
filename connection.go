package main

import (
  "github.com/gorilla/websocket"
  "log"
)

// Connection -
type Connection struct {
  ID      int
  ws      *websocket.Conn
  send    chan []byte
  receive *chan []byte
  closed  *chan int
}

// Listen - starts read and write loops
func (c *Connection) Listen() {
  go c.Writer()
  c.Reader()
  *c.closed <- c.ID
}

// Reader - reads message from websocket; puts in `receive` channel
func (c *Connection) Reader() {
  for {
    _, Message, err := c.ws.ReadMessage()
    if err != nil {
      log.Println("Couldn't Read:", err.Error())
      break
    }
    *c.receive <- Message
  }
  c.ws.Close()
}

// Writer - writes message to websocket
func (c *Connection) Writer() {
  for message := range c.send {
    err := c.ws.WriteMessage(websocket.TextMessage, message)
    if err != nil {
      log.Println("Couldn't Write:", err.Error())
      break
    }
  }
  c.ws.Close()
}

// Kill - closes channels
func (c *Connection) Kill() {
  close(c.send)
}
