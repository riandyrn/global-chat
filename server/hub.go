package main

import (
	"sync"
)

// Hub is used for bridging communication between sessions
type Hub struct {
	sessions  *sync.Map       // connected clients
	handles   *sync.Map       // list of active handles
	broadcast chan *MsgServer // channel for broadcasting messages to all connected clients
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

func (h *Hub) RegisterHandle(handle string) bool {
	_, isRegistered := h.handles.Load(handle)
	if !isRegistered {
		// register new handle
		h.handles.Store(handle, true)
	}
	return !isRegistered
}

func (h *Hub) run() {
	for {
		msg := <-h.broadcast
		// broadcast to all users
		h.sessions.Range(func(s, _ interface{}) bool {
			sess := s.(*Session)
			// skip skipHandle & unsigned session
			if sess.handle != msg.skipHandle && sess.handle != "" {
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
		handles:   &sync.Map{},
		broadcast: make(chan *MsgServer, 256),
	}
	go h.run()
	return h
}
