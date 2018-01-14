# Global Chat Room

Simple yet scalable global chat room application ;)

This project is meant to help me understand how to scale up chat application. Thus the feature of this application is very simple. User can only join to single global room, then user can publish the message where all of the users joined the room can read the message.

The server is written in golang & the client itself is written by using Vue.js & JQuery.

## Installation

1. Install [Go environment](https://golang.org/doc/install)
2. Go to `/server` directory, then execute build command: `go build`
3. Run server by using command: `./server`
4. Open your preferred browser then access `http://localhost:8192/`