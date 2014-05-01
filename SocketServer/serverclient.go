package SocketServer

import (
	"code.google.com/p/go.net/websocket"
	"log"
)

//Holds 1000 messages at once.
const channelSize = 1000

//holds all the info an client needs.
type Client struct {
	ws       *websocket.Conn
	server   *Server
	send     chan *MessageStruct
	done     chan bool
	username string
	userid   string
}

func NewClient(ws *websocket.Conn, server *Server, username, userid string) *Client {

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
	return &Client{ws, server, send, done, username, userid}
}

//getter for client connection
func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

func (c *Client) IP() string {
	return c.Conn().Request().RemoteAddr
}

//get write channel.
func (c *Client) Write() chan<- *MessageStruct {
	return c.send
}

func (c *Client) Listen() {
	go c.sendLoop()
	c.ListenToAll()
}

//return done channel.
func (c *Client) Done() chan<- bool {
	return c.done
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
			log.Println("Incoming message: " + message.Message + " from ip " + c.IP())
			c.server.BroadCast() <- &message
			break
		// START OMIT
		case "contact_list":
			contacts := c.server.GetContacts()
			usernames := make([]string, len(contacts))
			i := 0
			for contact, _ := range contacts {
				usernames[i] = contact.username
				i++
			}
			//struct containing all contacts.
			contactsMessage := &ContactMessage{
				Contacts: usernames,
				Type:     "contacts",
			}
			//send contact list to client
			websocket.JSON.Send(c.ws, &contactsMessage)
		//client requested a username
		// END OMIT
		case "user":
			ip := c.IP()
			userMessage := &MessageStruct{
				From:    ip,
				Message: ip,
				Type:    "user",
				Time:    message.Time,
			}
			//send to yourself!
			c.Write() <- userMessage
			break
		}
	}
}
