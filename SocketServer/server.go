package SocketServer

import (
	"code.google.com/p/go.net/websocket"
	"labix.org/v2/mgo"
	"log"
)

type Server struct {
	clients      map[*Client]bool
	addClient    chan *Client
	removeClient chan *Client
	sendAll      chan *MessageStruct
	messages     []*MessageStruct
	dbMessages   *mgo.Collection
	dbUsers      *mgo.Collection
}

func NewServer() *Server {
	server := Server{
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
	return (chan<- *Client)(s.addClient)
}

//channel to remove a client from the chat.
func (s *Server) RemoveClient() chan<- *Client {
	return (chan<- *Client)(s.removeClient)
}

//channel to broadcast the messages to all clients.
func (s *Server) BroadCast() chan<- *MessageStruct {
	return (chan<- *MessageStruct)(s.sendAll)
}

//holds all the messages from clients.
func (s *Server) Messages() []*MessageStruct {
	msgs := make([]*MessageStruct, len(s.messages))
	copy(msgs, s.messages)
	return msgs
}

func (s *Server) onConnectHandler() websocket.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		client := NewClient(ws, s)
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
			ip := newclient.getIP()
			log.Println("New client with ip " + ip + " added")
			s.clients[newclient] = true

			//write all previous messages to this client
			for _, msg := range s.messages {
				newclient.Write() <- msg
			}

		//client disconnected.
		case removeClient := <-s.removeClient:
			delete(s.clients, removeClient)
			log.Println("Client with ip " + removeClient.getIP() + " disconnected")

		//new message came in, distribute to all clients and db.
		case message := <-s.sendAll:
			s.messages = append(s.messages, message)
			for client, _ := range s.clients {
				client.Write() <- message
			}
		}
	}
}
