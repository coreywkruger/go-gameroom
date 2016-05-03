package main

import (
	"errors"
	"log"
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

		querystring := request.URL.Query()
		id := querystring["id"][0]

		ws, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			http.Error(writer, err.Error(), 500)
			return
		}

		conn := NewConnection(id, ConnConfig{receive: brc})

		// register new connection
		err = register(conn)
		if err != nil {
			http.Error(writer, err.Error(), 500)
			return
		}

		// read/write loops to websocket; read bocks until closed
		go conn.Writer(ws)
		conn.Reader(ws)

		// remove connection when closed
		_ = remove(conn.ID)
	}
}

// list of all members
var list = make(map[string]*Connection)

// register - creates new member
func register(c *Connection) error {
	if list[c.ID] != nil {
		return errors.New("id already in use")
	}
	list[c.ID] = c
	log.Println("registering:", c.ID, " # of connections: ", len(list))
	return nil
}

// GetAllConnections - gets all members
func GetAllConnections() (map[string]*Connection, error) {
	return list, nil
}

// remove - deletes a connection from the list
func remove(id string) error {
	delete(list, id)
	log.Println("removing:", id, " # of connections: ", len(list))
	return nil
}
