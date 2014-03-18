package SocketServer

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
	send   chan *MessageStruct
	done   chan bool
}

func NewClient(ws *websocket.Conn, server *Server) *Client {

	if ws == nil {
		log.Fatal("No connection")
	} else if server == nil {
		log.Fatal("No server")
	}

	//channel to send messages over.
	send := make(chan *MessageStruct, channelSize)

	//channel to prompt the server when done.
	done := make(chan bool)

	//returns new struct.
	return &Client{ws, server, send, done}
}

//getter for client connection
func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

func (c *Client) getIP() string {
	return c.Conn().Request().RemoteAddr
}

//get write channel.
func (c *Client) Write() chan<- *MessageStruct {
	return (chan<- *MessageStruct)(c.send)
}

func (c *Client) Listen() {
	go c.sendLoop()
	c.ListenToAll()
}

//return done channel.
func (c *Client) Done() chan<- bool {
	return (chan<- bool)(c.done)
}

func (c *Client) sendLoop() {
	for {
		select {
		case message := <-c.send:
			websocket.JSON.Send(c.ws, message)

		case <-c.done:
			c.server.RemoveClient() <- c
			c.done <- true
			return
		}
	}
}

func (c *Client) ListenToAll() {
	for {
		var message MessageStruct
		err := websocket.JSON.Receive(c.ws, &message)
		if err != nil {
			c.done <- true
			continue
		}
		//what message is coming?
		switch message.Type {
		case "message":
			log.Println("Incoming message: " + message.Message + " from ip " + c.getIP())
			c.server.BroadCast() <- &message
			break
		//client requested a username
		case "user":
			ip := c.getIP()
			userMessage := &MessageStruct{
				From:    ip,
				Message: ip,
				Type:    "user",
				Time:    message.Time,
			}
			c.Write() <- userMessage
			break
		}
	}
}
