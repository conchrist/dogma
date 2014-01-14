/* Copyright (C) 2013 Christopher Lillthors and Viktor Kronvall
 * This file is part of GoWebSocket.
 *
 * GoWebSocket is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * GoWebSocket is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with GoWebSocket.  If not, see <http://www.gnu.org/licenses/>.
 */

package SocketServer

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	pathToServer string
	clients      map[*Client]bool
	addClient    chan *Client
	removeClient chan *Client
	sendAll      chan *MessageStruct
	messages     []*MessageStruct
}

func NewServer(path string) *Server {
	tmp := Server{
		pathToServer: path,
		clients:      make(map[*Client]bool),
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
//why this??? copy all the messages to new client.
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
			log.Println("New client with ip " + newclient.getIP() + " added")
			s.clients[newclient] = true

			//write all previous messages to clients
			for _, msg := range s.messages {
				newclient.Write() <- msg
			}

		//client disconnected.
		case removeClient := <-s.removeClient:
			delete(s.clients, removeClient)
			log.Println("Client with ip " + removeClient.getIP() + " disconnected")

		case sendall := <-s.sendAll:
			message := sendall
			//fmt.Println("Broadcast all the messages")
			s.messages = append(s.messages, message)
			for client, _ := range s.clients {
				client.Write() <- message
			}
		}
	}
}
