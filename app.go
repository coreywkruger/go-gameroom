package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// App -
type App struct{}

// NewApp -
func NewApp() App {
	return App{}
}

// Listen -
func (a *App) Listen(port string) {

	cr := ConnectionRegistry{}

	// Main http rounter
	r := mux.NewRouter()
	r.Headers("Access-Control-Allow-Origin", "*")
	r.Headers("Access-Control-Allow-Headers", "Content-Type")
	r.Headers("Content-Type", "application/json")
	r.Headers("Content-Type", "text/plain")
	r.HandleFunc("/register", RegistrationHandler(&cr))

	// listens for i/o from ws
	cr.StartBroadcast()

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
