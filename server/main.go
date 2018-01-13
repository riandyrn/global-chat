package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const listenPort = 8192

var hub *Hub

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	// initialize hub for communicating between session
	hub = NewHub()

	// handle main page
	http.Handle("/", http.FileServer(http.Dir("../client")))

	// handle websocket connections
	http.HandleFunc("/wsc", handleWebsocketConn)

	// start server
	log.Printf("http server started on :%d", listenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil))
}

func handleWebsocketConn(w http.ResponseWriter, r *http.Request) {
	// attempt to upgrade http connection to websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("unable to upgrade connection for: %v, due: %v", r.RemoteAddr, err)
		return
	}
	defer ws.Close()

	sess := NewSession(ws, hub)
	hub.attachSession(sess)
	sess.ReadUserCommands() // will block until error
}
