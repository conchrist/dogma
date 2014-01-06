package Chat

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"net/http"
	"strconv"
)

type Server struct {
	pathToServer string
	clients      []*Client
	addClient    chan *Client
	removeClient chan *Client
	sendAll      chan *MessageStruct
	messages     []*MessageStruct
}

func NewServer(path string) *Server {
	tmp := Server{
		pathToServer: path,
		clients:      make([]*Client, 0),
		addClient:    make(chan *Client),
		removeClient: make(chan *Client),
		sendAll:      make(chan *MessageStruct),
		messages:     make([]*MessageStruct, 0),
	}
	return &tmp
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
//why this???
func (s *Server) Messages() []*MessageStruct {
	msgs := make([]*MessageStruct, len(s.messages))
	copy(msgs, s.messages)
	return msgs
}

//start server!
func (s *Server) Listen() {
	fmt.Println("Server listening")

	onConnect := func(ws *websocket.Conn) {
		client := NewClient(ws, s)
		s.addClient <- client
		//not yet implemented...
		client.Listen()
		defer ws.Close()
	}
	//new connections will have this handler.
	http.Handle(s.pathToServer, websocket.Handler(onConnect))

	//listen for messages, clients and so on...

	for {
		select {

		//new client connecting
		case newclient := <-s.addClient:
			fmt.Println("New client added")
			s.clients = append(s.clients, newclient)
			//write all the messages to client
			for _, msg := range s.messages {
				newclient.Write() <- msg
			}
			fmt.Println("Connected clients: " + strconv.Itoa(len(s.clients)))

		//client disconnected.
		case removeClient := <-s.removeClient:
			fmt.Println("Remove client")
			for i := range s.clients { //i is index for the client to remove
				if s.clients[i] == removeClient {
					//remove client from slice.
					s.clients = append(s.clients[:i], s.clients[i+1:]...)
					break //come out of range loop.
				}
			} //end of for

		case sendall := <-s.sendAll:
			fmt.Println("Broadcast all the messages")
			s.messages = append(s.messages, sendall)
			for _, c := range s.clients {
				c.Write() <- sendall
			}
		}
	}
} //end of listen()
