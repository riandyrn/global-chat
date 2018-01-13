package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// Session is used to represent user session
type Session struct {
	handle string
	conn   *websocket.Conn
	hub    *Hub
}

// ReadIncomingMessages is used for processing user input
// will block until error while reading input
func (s *Session) ReadIncomingMessages() {
	for {
		// read client message
		var msg MsgClient
		if err := s.conn.ReadJSON(&msg); err != nil {
			log.Printf("unable to read incoming message for: %v, due: %v", s.handle, err)
			s.hub.detachSession(s)
			return
		}
		// parse client message
		switch {
		case msg.Join != nil:
			s.join(&msg)
		case msg.Pub != nil:
			s.publish(&msg)
		}
	}
}

func (s *Session) join(msg *MsgClient) {
	s.handle = msg.Join.Handle
	resp := &MsgServer{
		What: "join",
		From: s.handle,
	}
	s.hub.broadcast <- resp
}

func (s *Session) publish(msg *MsgClient) {
	resp := &MsgServer{
		What:    "msg",
		From:    s.handle,
		Content: msg.Pub.Content,
	}
	s.hub.broadcast <- resp
}

// Destroy is used for destroying session
func (s *Session) Destroy() {
	resp := &MsgServer{
		What: "left",
		From: s.handle,
	}
	s.hub.broadcast <- resp
	s.conn.Close()
}

// NewSession is used for creating new session instance
func NewSession(conn *websocket.Conn, hub *Hub) *Session {
	sess := &Session{conn: conn, hub: hub}
	return sess
}
