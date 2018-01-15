package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var hub *Hub
var debugMode *bool

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func logDebugMessage(format string, args ...interface{}) {
	if !*debugMode {
		return
	}
	if len(args) > 0 {
		log.Printf(format, args...)
		return
	}
	log.Println(format)
}

func main() {

	clientPath := flag.String("client", "../client", "Path to client files")
	listenPort := flag.Int("port", 8192, "Default port to listen")
	debugMode = flag.Bool("debug", false, "Turn on debug mode?")
	flag.Parse()

	// initialize hub for communicating between session
	hub = NewHub()

	// handle main page
	http.Handle("/", http.FileServer(http.Dir(*clientPath)))

	// handle websocket connections
	http.HandleFunc("/wsc", handleWebsocketConn)

	// start server
	log.Printf("http server started on :%d", *listenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *listenPort), nil))
}

func handleWebsocketConn(w http.ResponseWriter, r *http.Request) {
	// attempt to upgrade http connection to websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logDebugMessage("unable to upgrade connection for: %v, due: %v", r.RemoteAddr, err)
		return
	}
	defer ws.Close()
	logDebugMessage("new socket connection initiated")

	sess := NewSession(ws, hub)
	hub.attachSession(sess)
	sess.ReadUserCommands() // will block until error
}
