package SocketServer

import (
	"code.google.com/p/go.net/websocket"
	"labix.org/v2/mgo"
	"log"
	"sync"
)

type Server struct {
	mutex        *sync.Mutex
	clients      map[*Client]bool
	addClient    chan *Client
	removeClient chan *Client
	sendAll      chan *MessageStruct
	messages     []*MessageStruct
}

func NewServer(address, name string) *Server {
	server := Server{
		mutex:        &sync.Mutex{},
		clients:      make(map[*Client]bool),
		addClient:    make(chan *Client),
		removeClient: make(chan *Client),
		sendAll:      make(chan *MessageStruct),
		messages:     make([]*MessageStruct, 0),
	}
	return &server
}

//channel to add a client to the chat.
func (s *Server) AddClient() chan<- *Client {
	return (s.addClient)
}

//channel to remove a client from the chat.
func (s *Server) RemoveClient() chan<- *Client {
	return (s.removeClient)
}

//channel to broadcast the messages to all clients.
func (s *Server) BroadCast() chan<- *MessageStruct {
	return (s.sendAll)
}

//holds all the messages from clients.
func (s *Server) Messages() []*MessageStruct {
	msgs := make([]*MessageStruct, len(s.messages))
	copy(msgs, s.messages)
	return msgs
}

//return the contact list. (A list of users)
func (s *Server) GetContacts() map[*Client]bool {
	return s.clients
}

//a handler to handle a new client connection.
func (s *Server) onConnectHandler(username, userid string, db *mgo.Database) websocket.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		client := NewClient(ws, s, username, userid, db)
		s.addClient <- client
		client.Listen()
	})
}

//start server!
func (s *Server) Listen() {
	//listen for messages, clients and so on...
	for {
		select {
		//new client connecting
		case newclient := <-s.addClient:
			ip := newclient.IP()
			log.Println("New client with ip " + ip + " added")
			s.mutex.Lock()
			s.clients[newclient] = true
			s.mutex.Unlock()

			//write all previous messages to this client
			for _, msg := range s.messages {
				newclient.Write() <- msg
			}

		//client disconnected.
		case removeClient := <-s.removeClient:
			s.mutex.Lock()
			delete(s.clients, removeClient)
			s.mutex.Unlock()
			log.Println("Client with ip " + removeClient.IP() + " disconnected")

		//new message came in, distribute to all clients.
		case message := <-s.sendAll:
			s.messages = append(s.messages, message)
			for client, _ := range s.clients {
				client.Write() <- message
			}
		}
	}
}
