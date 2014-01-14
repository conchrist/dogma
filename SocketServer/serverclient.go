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
	"log"
	"time"
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

//get write channel. Implements Write method! :)
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
	//log.Println("Write to all")

	for {
		select {
		case message := <-c.send:
			//log.Println("Sending... ", message)
			websocket.JSON.Send(c.ws, message)

		case <-c.done:
			//log.Println("Remove yourself!")
			c.server.RemoveClient() <- c
			c.done <- true
			return
		}
	}
}

func (c *Client) ListenToAll() {
	//log.Println("Listening to all clients")
	for {
		select {
		case <-c.done:
			//log.Println("Remove yourself!")
			c.server.RemoveClient() <- c
			c.done <- true
			return
		default:
			var message MessageStruct
			err := websocket.JSON.Receive(c.ws, &message)
			if err != nil {
				//something is wrong, close yourself
				c.done <- true
				continue
			}
			switch message.Type {
			case "message":
				log.Printf("Message recieved %s\n", message.Message)
				c.server.BroadCast() <- &message
				break
			case "user":
				log.Printf("Username requested by %s\n", c.getIP())
				userMessage := &MessageStruct{
					From:    c.getIP(),
					Message: c.getIP(),
					Type:    "user",
					Time:    int(time.Now().UnixNano() % 1e6 / 1e3),
				}
				c.Write() <- userMessage
				break
			}
		}
	}
}
