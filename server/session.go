package main

import (
	"encoding/json"
	"log"
	"time"

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
		now := timeNow()
		if _, raw, err := s.conn.ReadMessage(); err != nil {
			log.Printf("unable to read incoming message for: %v, due: %v", s.handle, err)
			s.hub.detachSession(s)
			return
		} else if err = json.Unmarshal(raw, &msg); err != nil {
			s.QueueOut(ErrMalformed("", now))
			continue
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
	id := msg.Join.ID

	// check already join error
	if s.handle != "" {
		s.QueueOut(ErrAlreadyJoin(id, now))
		return
	}
	// check bad request error
	if msg.Join.Handle == "" {
		s.QueueOut(ErrMalformed(id, now))
		return
	}
	// check whether handle is taken already by other user
	if s.hub.isHandleTaken(msg.Join.Handle) {
		s.QueueOut(ErrHandleTaken(id, now))
		return
	}

	// register handle
	s.handle = msg.Join.Handle
	s.hub.regHandle <- s.handle

	// output success to user
	s.QueueOut(NoErr(id, "join", now))

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
	id := msg.Pub.ID

	// check command out of sequence error
	if s.handle == "" {
		s.QueueOut(ErrCommandOutOfSequence(id, now))
		return
	}
	// check bad request error
	if msg.Pub.Content == "" {
		s.QueueOut(ErrMalformed(id, now))
		return
	}

	// notify user message has been accepted
	s.QueueOut(NoErrAccepted(id, "pub", now))

	// notify other user
	s.broadcastMessage(&MsgServer{
		Data: &DataPayload{
			From:      s.handle,
			Content:   msg.Pub.Content,
			Timestamp: now,
		},
	})
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

func (s *Session) broadcastMessage(msg *MsgServer) {
	select {
	case s.hub.broadcast <- msg:
	default:
		log.Printf("hub queue is full, dropping message for: %v", s.handle)
	}
}

// QueueOut is used for buffering server response
func (s *Session) QueueOut(msg *MsgServer) {
	select {
	case s.respQueue <- msg:
	case <-time.After(50 * time.Microsecond):
		log.Printf("QueueOut timeout for: %v", s.handle)
	}
}

// Destroy is used for destroying session
func (s *Session) Destroy() {
	now := timeNow()
	// notify other users, user is leaving
	s.broadcastMessage(&MsgServer{
		Pres: &PresPayload{
			What:      "left",
			From:      s.handle,
			Timestamp: now,
		},
		skipHandle: s.handle, // skip user session
	})
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
