package main

// MsgClient represents client input
type MsgClient struct {
	Join *JoinPayload `json:"join,omitempty"`
	Pub  *PubPayload  `json:"pub,omitempty"`
}

// JoinPayload represents join command payload
type JoinPayload struct {
	Handle string `json:"handle"`
}

// PubPayload represents publish command payload
type PubPayload struct {
	Content string `json:"content"`
}

// MsgServer represents server response
type MsgServer struct {
	// message context: msg, join, left
	What string `json:"what,omitempty"`
	// message origin
	From string `json:"from,omitempty"`
	// message content
	Content string `json:"content,omitempty"`
}
