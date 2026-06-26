package main

import (
	"fmt"
	"net/http"
	"study/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/echo", handlers.HandleSendEcho).Methods("GET", "OPTIONS")
	r.HandleFunc("/telemetry", handlers.HandleSendTelemetry).Methods("GET", "OPTIONS")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	http.Handle("/", r)

	fmt.Println("Server has been started successfully!")
	http.ListenAndServe(":8080", nil)
}
