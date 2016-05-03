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
func RegistrationHandler(writer http.ResponseWriter, request *http.Request) {

	// websocket stuff
	ws, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	err = registerConnection(request.URL.Query()["id"][0], ws)
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}
}

func registerConnection(id string, ws *websocket.Conn) error {
	// create new connection struct
	conn := NewConnection(id, ConnConfig{
		receive: &broadcastChannel,
	})

	// register new connection
	err := storeConnection(conn)
	if err != nil {
		return err
	}
	// read/write loops to websocket; read bocks until closed
	conn.Start(ws)
	// remove connection when closed
	return remove(conn.ID)
}

// StartBroadcast -
func StartBroadcast() {
	// listens for i/o from ws
	go func() {
		for {
			select {
			case broadcast := <-broadcastChannel:
				// broadcastChannel received; broadcast to all connections
				connections, _ := GetAllConnections()
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

// ConnectionList of all members
var ConnectionList = make(map[string]*Connection)

// register - creates new member
func storeConnection(c *Connection) error {
	if ConnectionList[c.ID] != nil {
		return errors.New("id already in use")
	}
	ConnectionList[c.ID] = c
	log.Println("registering:", c.ID, " # of connections: ", len(ConnectionList))
	return nil
}

// GetAllConnections - gets all members
func GetAllConnections() (map[string]*Connection, error) {
	return ConnectionList, nil
}

// remove - deletes a connection from the ConnectionList
func remove(id string) error {
	delete(ConnectionList, id)
	log.Println("removing:", id, " # of connections: ", len(ConnectionList))
	return nil
}
