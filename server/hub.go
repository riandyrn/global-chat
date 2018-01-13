package main

import (
	"sync"
)

// Hub is used for bridging communication between sessions
type Hub struct {
	sessions  *sync.Map       // connected clients
	broadcast chan *MsgServer // channel for broadcasting messages to all connected clients
}

func (h *Hub) attachSession(sess *Session) {
	h.sessions.Store(sess, true)
}

func (h *Hub) detachSession(sess *Session) {
	sess.Destroy()
	h.sessions.Delete(sess)
}

func (h *Hub) run() {
	for {
		msg := <-h.broadcast
		h.sessions.Range(func(s, _ interface{}) bool {
			sess := s.(*Session)
			if sess.handle != msg.skipHandle {
				sess.QueueOut(msg)
			}
			return true
		})
	}
}

// NewHub is used for creating new hub instance
func NewHub() *Hub {
	h := &Hub{
		sessions:  &sync.Map{},
		broadcast: make(chan *MsgServer, 4096),
	}
	go h.run()
	return h
}
