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
func RegistrationHandler(conn *Connection) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		ws, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			http.Error(writer, err.Error(), 500)
			return
		}

		_ = RegisterConnection(conn)

		go conn.Writer(ws)
		conn.Reader(ws)
		*conn.closed <- conn.ID
	}
}

// list of all members
var list = make(map[int]*Connection)

// RegisterConnection - creates new member
func RegisterConnection(c *Connection) error {
	log.Println("registering ", c.ID)
	list[c.ID] = c
	log.Println("# of members: ", len(list))
	return nil
}

// GetConnection - gets member by id
func GetConnection(id int) (*Connection, error) {
	m := list[id]
	if m != nil {
		log.Println("Getting Connection: ", m.ID)
		return m, nil
	}
	return nil, errors.New("Connection not found")
}

// GetAllConnections - gets all members
func GetAllConnections() (map[int]*Connection, error) {
	return list, nil
}
