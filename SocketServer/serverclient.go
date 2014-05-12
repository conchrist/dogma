package SocketServer

import (
	"code.google.com/p/go.net/websocket"
	"errors"
	"labix.org/v2/mgo"
	"log"
)

var (
	DBERROR        error = errors.New("Error dialing DB")
	DBINSERT       error = errors.New("Error inserting value")
	WEBSOCKETERROR error = errors.New("Client failed to connect")
	SERVERERROR    error = errors.New("Server not started")
)

//Holds 1000 messages at once.
const channelSize = 1000

//holds all the info an client needs.
type Client struct {
	ws       *websocket.Conn
	server   *server
	send     chan *messageStruct
	done     chan bool
	username string
	userid   string
	db       *mgo.Database
}

func NewClient(ws *websocket.Conn, server *server, username, userid string, db *mgo.Database) *Client {

	if ws == nil {
		log.Fatal(WEBSOCKETERROR.Error())
	} else if server == nil {
		log.Fatal(SERVERERROR.Error())
	}

	//channel to send messages over.
	send := make(chan *messageStruct, channelSize)

	//channel to prompt the server when done.
	done := make(chan bool)

	//returns new struct.
	return &Client{
		ws:       ws,
		server:   server,
		send:     send,
		done:     done,
		username: username,
		userid:   userid,
		db:       db,
	}
}

//insert messages into DB.
func (c *Client) insertMessage(message *messageStruct) {
	err := c.db.C("Messages").Insert(message)
	if err != nil {
		log.Fatalf("%s %s", DBINSERT.Error(), err.Error())
	}
}

//getter for client connection
func (c *Client) conn() *websocket.Conn {
	return c.ws
}

func (c *Client) iP() string {
	return c.Conn().Request().RemoteAddr
}

//get write channel.
func (c *Client) write() chan<- *messageStruct {
	return c.send
}

//start client and listen on connections.
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
		var message messageStruct
		err := websocket.JSON.Receive(c.ws, &message)
		if err != nil {
			c.done <- true
			continue
		}
		//what message is coming?
		switch message.Type {
		case "message":
			log.Println("Incoming message: " + message.Message + " from ip " + c.IP())
			go c.insertMessage(&message)
			c.server.BroadCast() <- &message
			break
		// START OMIT
		case "contact_list":
			contacts := c.server.GetContacts()

			c.server.mutex.Lock()
			usernames := make([]string, len(contacts))
			i := 0
			for contact, _ := range contacts {
				usernames[i] = contact.username
				i++
			}
			c.server.mutex.Unlock()
			//struct containing all contacts.
			contactsMessage := &contactMessage{
				Contacts: usernames,
				Type:     "contacts",
			}
			//send contact list to client
			websocket.JSON.Send(c.ws, &contactsMessage)
		//client requested a username
		// END OMIT
		case "user":
			ip := c.IP()
			userMessage := &messageStruct{
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
