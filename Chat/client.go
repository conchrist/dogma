package Chat

import (
	"code.google.com/p/go.net/websocket"
	"log"
)

//Holds 1000 messages at once.
const channelSize = 1000

//holds all the info an client needs.
type Client struct {
	ws     *websocket.Conn
	server *Server
	ch     chan *MessageStruct
	done   chan bool
}

func NewClient(ws *websocket.Conn, server *Server) *Client {

	if ws == nil {
		log.Fatal("No connection")
	} else if server == nil {
		log.Fatal("No server")
	}

	//channel to send messages over.
	ch := make(chan *MessageStruct, channelSize)

	//channel to prompt the server when done.
	done := make(chan bool)

	//returns new struct.
	return &Client{ws, server, ch, done}
}

//getter for client connection
func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

//get write channel. Implements Write method! :)
func (c *Client) Write() chan<- *MessageStruct {
	return (chan<- *MessageStruct)(c.ch)
}

//return done channel.
func (c *Client) Done() chan<- bool {
	return (chan<- bool)(c.done)
}

func (c *Client) writeToAll() {
	log.Println("Write to all")

	for {
		select {
		case message := <-c.ch:

		}
	}
}
