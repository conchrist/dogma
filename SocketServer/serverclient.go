package SocketServer

import (
	"code.google.com/p/go.net/websocket"
	"errors"
	"github.com/wsxiaoys/terminal/color"
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
type client struct {
	ws       *websocket.Conn
	server   *server
	send     chan *messageStruct
	_done    chan bool
	username string
	userid   string
	db       *mgo.Database
	running  bool
}

func NewClient(ws *websocket.Conn, server *server, username, userid string, db *mgo.Database) *client {

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
	return &client{
		ws:       ws,
		server:   server,
		send:     send,
		_done:    done,
		username: username,
		userid:   userid,
		db:       db,
		running:  false,
	}
}

//insert messages into DB.
func (c *client) insertMessage(message *messageStruct) {
	err := c.db.C("Messages").Insert(message)
	if err != nil {
		log.Fatalf("%s %s", DBINSERT.Error(), err.Error())
	}
}

//getter for client connection
func (c *client) conn() *websocket.Conn {
	return c.ws
}

func (c *client) iP() string {
	return c.conn().Request().RemoteAddr
}

//get write channel.
func (c *client) write() chan<- *messageStruct {
	return c.send
}

//start client and listen on connections.
func (c *client) Listen() {
	c.running = true
	go c.sendLoop()
	c.listenToAll()
}

//return done channel.
func (c *client) done() chan<- bool {
	return c._done
}

func (c *client) sendLoop() {
	for c.running {
		select {
		case message := <-c.send:
			websocket.JSON.Send(c.ws, message)

		case <-c._done:
			c.server.RemoveClient() <- c
			c._done <- true
			return
		}
	}
}

func (c *client) listenToAll() {
	for c.running {
		var message messageStruct
		err := websocket.JSON.Receive(c.ws, &message)
		if err != nil {
			c._done <- true
			continue
		}
		//what message is coming?
		switch message.Type {
		case "message":
			color.Printf("@{mK}Incoming message: %s from ip %s", message.Message, c.iP())
			color.Println()
			c.server.BroadCast() <- &message
			go c.insertMessage(&message)
			break
		// START OMIT
		case "image":
			log.Println("Image received from client " + message.From + "from ip " + c.iP())
			c.server.BroadCast() <- &message
			go c.insertMessage(&message)
		case "contact_list":
			contacts := c.server.getContacts()

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
			ip := c.iP()
			userMessage := &messageStruct{
				From:    ip,
				Message: ip,
				Type:    "user",
				Time:    message.Time,
			}
			//send to yourself!
			c.write() <- userMessage
			break
		}
	}
}

func (c *client) Close() {
	c.running = false
}
