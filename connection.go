package main

import (
  "github.com/gorilla/websocket"
  "log"
)

// Connection contains websocket connnetion
type Connection struct {
  ID       int
  ws       *websocket.Conn
  outbound chan []byte
}

// Init -
func (c *Connection) Init(broadcast chan []byte) {
  // defer func() {
  //   close(c.outbound)
  //   broadcast <- []byte(strconv.Itoa(c.ID))
  // }()
  go c.Writer()
  c.Reader(broadcast)
}

// Reader -
func (c *Connection) Reader(input chan []byte) {
  for {
    _, Bytes, err := c.ws.ReadMessage()

    if err != nil {
      log.Println(err.Error())
      break
    }
    input <- Bytes
  }
  c.ws.Close()
}

// Writer -
func (c *Connection) Writer() {
  for message := range c.outbound {
    log.Println("broadcasting this: ", message)
    err := c.ws.WriteMessage(websocket.TextMessage, message)
    if err != nil {
      log.Println(err.Error())
      break
    }
  }
  c.ws.Close()
}
