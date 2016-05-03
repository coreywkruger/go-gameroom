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

var broadcastChannel = make(chan Message)

// RegistrationHandler - registers new member & connection
func RegistrationHandler(cr *ConnectionRegistry) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// websocket stuff
		ws, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			http.Error(writer, err.Error(), 500)
			return
		}

		// create new connection struct
		conn := NewConnection(request.URL.Query()["id"][0], ConnConfig{
			receive: &broadcastChannel,
		})

		// register new connection
		err = cr.Save(conn)
		if err != nil {
			http.Error(writer, err.Error(), 500)
			return
		}
		// read/write loops to websocket; read bocks until closed
		conn.Start(ws)
		// Remove connection when closed
		err = cr.Remove(conn.ID)
		if err != nil {
			http.Error(writer, err.Error(), 500)
			return
		}
	}
}

// ConnectionRegistry -
type ConnectionRegistry struct {
	list map[string]*Connection
}

// StartBroadcast -
func (cr *ConnectionRegistry) StartBroadcast() {
	// listens for i/o from ws
	go func() {
		for {
			select {
			case broadcast := <-broadcastChannel:
				// broadcastChannel received; broadcast to all connections
				connections, _ := cr.GetAllConnections()
				log.Println("Broadcasting", string(broadcast.message))
				for _, connection := range connections {
					if connection.ID != broadcast.ID {
						*connection.send <- broadcast.message
					}
				}
			}
		}
	}()
}

// Save - creates new member
func (cr *ConnectionRegistry) Save(c *Connection) error {
	if cr.list[c.ID] != nil {
		return errors.New("id already in use")
	}
	cr.list[c.ID] = c
	log.Println("registering:", c.ID, " # of connections: ", len(cr.list))
	return nil
}

// GetAllConnections - gets all members
func (cr *ConnectionRegistry) GetAllConnections() (map[string]*Connection, error) {
	return cr.list, nil
}

// Remove - deletes a connection from the cr.list
func (cr *ConnectionRegistry) Remove(id string) error {
	delete(cr.list, id)
	log.Println("removing:", id, " # of connections: ", len(cr.list))
	return nil
}
