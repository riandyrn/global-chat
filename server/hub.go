package main

import (
	"sync"
)

// Hub is used for bridging communication between sessions
type Hub struct {
	sessions  *sync.Map       // connected clients
	handles   *sync.Map       // list of active handles
	broadcast chan *MsgServer // channel for broadcasting messages to all connected clients
	regHandle chan string     // channel for registering handle
}

func (h *Hub) attachSession(sess *Session) {
	h.sessions.Store(sess, true)
}

func (h *Hub) detachSession(sess *Session) {
	sess.Destroy()
	h.sessions.Delete(sess)
	if sess.handle != "" {
		h.handles.Delete(sess.handle)
	}
}

func (h *Hub) isHandleTaken(handle string) bool {
	_, ok := h.handles.Load(handle)
	return ok
}

func (h *Hub) run() {
	for {
		select {
		case msg := <-h.broadcast:
			// broadcast to all users
			h.sessions.Range(func(s, _ interface{}) bool {
				sess := s.(*Session)
				// skip skipHandle & unsigned session
				if sess.handle != msg.skipHandle && sess.handle != "" {
					sess.QueueOut(msg)
				}
				return true
			})
		case handle := <-h.regHandle:
			// register new handle
			h.handles.Store(handle, true)
		}
	}
}

// NewHub is used for creating new hub instance
func NewHub() *Hub {
	h := &Hub{
		sessions:  &sync.Map{},
		handles:   &sync.Map{},
		broadcast: make(chan *MsgServer, 256),
		regHandle: make(chan string, 256),
	}
	go h.run()
	return h
}
