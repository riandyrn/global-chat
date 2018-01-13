package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// Session is used to represent user session
type Session struct {
	handle    string
	conn      *websocket.Conn
	hub       *Hub
	respQueue chan *MsgServer
}

// ReadUserCommands is used for processing user input
// will block until error while reading input
func (s *Session) ReadUserCommands() {
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
	now := timeNow()

	// check already join error
	if s.handle != "" {
		s.QueueOut(ErrAlreadyJoin(msg.Join.ID, now))
		return
	}
	// check bad request error
	if msg.Join.Handle == "" {
		s.QueueOut(ErrBadRequest(msg.Join.ID, now))
		return
	}

	// TODO: check whether handle is taken already by other user
	s.handle = msg.Join.Handle
	// output success to user
	s.QueueOut(NoErr(msg.Join.ID, "join", now))

	// notify other users
	s.hub.broadcast <- &MsgServer{
		Pres: &PresPayload{
			What:      "join",
			From:      s.handle,
			Timestamp: now,
		},
		skipHandle: s.handle, // skip user session
	}
}

func (s *Session) publish(msg *MsgClient) {
	now := timeNow()

	// check command out of sequence error
	if s.handle == "" {
		s.QueueOut(ErrCommandOutOfSequence(msg.Pub.ID, now))
		return
	}
	// check bad request error
	if msg.Pub.Content == "" {
		s.QueueOut(ErrBadRequest(msg.Pub.ID, now))
		return
	}

	// notify other user
	s.hub.broadcast <- &MsgServer{
		Data: &DataPayload{
			From:      s.handle,
			Content:   msg.Pub.Content,
			Timestamp: now,
		},
	}
}

func (s *Session) consumeQueue() {
	for {
		msg := <-s.respQueue
		if err := s.conn.WriteJSON(msg); err != nil {
			log.Printf("unable to broadcast message to: %v, due: %v", s.handle, err)
			s.hub.detachSession(s)
		}
	}
}

// QueueOut is used for buffering server response
func (s *Session) QueueOut(msg *MsgServer) {
	s.respQueue <- msg
}

// Destroy is used for destroying session
func (s *Session) Destroy() {
	now := timeNow()
	resp := &MsgServer{
		Pres: &PresPayload{
			What:      "left",
			From:      s.handle,
			Timestamp: now,
		},
		skipHandle: s.handle, // skip user session
	}
	s.hub.broadcast <- resp
	s.conn.Close()
}

// NewSession is used for creating new session instance
func NewSession(conn *websocket.Conn, hub *Hub) *Session {
	sess := &Session{
		conn:      conn,
		hub:       hub,
		respQueue: make(chan *MsgServer, 256),
	}
	go sess.consumeQueue()
	return sess
}
