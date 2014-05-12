package SocketServer

import (
	"github.com/codegangsta/martini"
	"labix.org/v2/mgo"
	"log"
)

//specify how a message looks like.
type messageStruct struct {
	From    string `json:"from"`
	Message string `json:"body"`
	Type    string `json:"type"`
	Time    int    `json:"time"`
}

//specify how the contact list looks like.
type contactMessage struct {
	Contacts []string `json:"contacts"`
	Type     string   `json:"type"`
}

//middleware
func mongoDB(adress, db string) martini.Handler {
	session, err := mgo.Dial(adress)
	if err != nil {
		log.Fatal("Could not connect to database " + err.Error())
	}
	return func(c martini.Context) {
		s := session.Clone()
		c.Map(s.DB(db))
		defer s.Close()
		c.Next()
	}
}
