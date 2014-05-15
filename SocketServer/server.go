package SocketServer

import (
	"code.google.com/p/go.net/websocket"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"sync"
)

var running = true

type server struct {
	mutex        *sync.Mutex
	clients      map[*client]bool
	addClient    chan *client
	removeClient chan *client
	sendAll      chan *messageStruct
	_messages    []*messageStruct
	running      bool
}

func NewServer(address, name string) *server {
	session, err := mgo.Dial(address)
	defer session.Close()

	messages := make([]*messageStruct, 0)

	err = session.DB(name).C("Messages").Find(bson.M{}).All(&messages)
	if err != nil {
		log.Fatalf("%s %s", DBERROR.Error(), err.Error())
		return nil
	}

	server := &server{
		mutex:        &sync.Mutex{},
		clients:      make(map[*client]bool),
		addClient:    make(chan *client),
		removeClient: make(chan *client),
		sendAll:      make(chan *messageStruct),
		_messages:    messages,
		running:      false,
	}
	return server
}

//channel to add a client to the chat.
func (s *server) AddClient() chan<- *client {
	return s.addClient
}

//channel to remove a client from the chat.
func (s *server) RemoveClient() chan<- *client {
	return s.removeClient
}

//channel to broadcast the messages to all clients.
func (s *server) BroadCast() chan<- *messageStruct {
	return s.sendAll
}

//holds all the messages from clients.
func (s *server) messages() []*messageStruct {
	msgs := make([]*messageStruct, len(s._messages))
	copy(msgs, s._messages)
	return msgs
}

//return the contact list. (A list of users)
func (s *server) getContacts() map[*client]bool {
	return s.clients
}

//a handler to handle a new client connection.
func (s *server) onConnectHandler(username, userid string, db *mgo.Database) websocket.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		client := NewClient(ws, s, username, userid, db)
		s.addClient <- client
		client.Listen()
	})
}

//start server!
func (s *server) Listen() {
	s.running = true
	//listen for messages, clients and so on...
	for s.running {
		select {
		//new client connecting
		case newclient := <-s.addClient:
			ip := newclient.iP()
			log.Println("New client with ip " + ip + " added")
			s.mutex.Lock()
			s.clients[newclient] = true
			s.mutex.Unlock()

			//write all previous messages to this client
			for _, msg := range s._messages {
				newclient.write() <- msg
			}

			go func() {
				s.BroadCast() <- &messageStruct{
					From:    "server",
					Message: newclient.username,
					Type:    "client joined",
				}
			}()

		//client disconnected.
		case removeClient := <-s.removeClient:
			s.mutex.Lock()
			delete(s.clients, removeClient)
			s.mutex.Unlock()
			log.Println("Client with ip " + removeClient.iP() + " disconnected")
			go func() {
				s.BroadCast() <- &messageStruct{
					From:    "server",
					Message: removeClient.username,
					Type:    "client left",
				}
			}()

		//new message came in, distribute to all clients.
		case message := <-s.sendAll:
			s._messages = append(s._messages, message)
			for client, _ := range s.clients {
				client.write() <- message
			}
		}
	}
}

func (s *server) Close() {
	s.running = false
	for client, _ := range s.clients {
		client.Close()
	}
}
