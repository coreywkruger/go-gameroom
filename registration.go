package main

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// RegistrationHandler - registers new member & connection
func RegistrationHandler(brc *chan Message) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// make new connection
		send := make(chan []byte, 256)
		conn := &Connection{
			ID:      rand.Intn(1000),
			send:    &send,
			receive: brc,
		}

		ws, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			http.Error(writer, err.Error(), 500)
			return
		}

		// register new connection
		_ = register(conn)

		// read/write loops to websocket; read bocks until closed
		go conn.Writer(ws)
		conn.Reader(ws)

		// remove connection when closed
		remove(conn.ID)
	}
}

// list of all members
var list = make(map[int]*Connection)

// register - creates new member
func register(c *Connection) error {
	list[c.ID] = c
	log.Println("registering:", c.ID, " # of connections: ", len(list))
	return nil
}

// GetAllConnections - gets all members
func GetAllConnections() (map[int]*Connection, error) {
	return list, nil
}

// remove - deletes a connection from the list
func remove(id int) error {
	delete(list, id)
	log.Println("removing:", id, " # of connections: ", len(list))
	return nil
}
