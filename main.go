package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var broadcastChannel = make(chan []byte)
var closingChannel = make(chan int)

func main() {

	// Main http rounter
	r := mux.NewRouter()
	r.Headers("Access-Control-Allow-Origin", "*")
	r.Headers("Access-Control-Allow-Headers", "Content-Type")
	r.Headers("Content-Type", "application/json")
	r.Headers("Content-Type", "text/plain")
	r.HandleFunc("/register", CreateRegistrationHandler(broadcastChannel, closingChannel))

	connections := GetAllConnections()
	// listens for i/o from ws
	go func() {
		for {
			select {
			case broadcast := <-broadcastChannel:
				// broadcastChannel received; broadcast to all connections
				log.Println("Broadcasting", string(broadcast))
				for _, connection := range connections {
					connection.send <- broadcast
				}
			case id := <-closingChannel:
				// connection is signaling close; kill & delete connection
				log.Println("Closing", id)
				connections[id].Kill()
				delete(connections, id)
			}
		}
	}()

	err := http.ListenAndServe(":3334", r)

	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
