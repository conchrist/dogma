package SocketServer

import (
	"github.com/codegangsta/martini"
	"labix.org/v2/mgo"
)

type MessageStruct struct {
	From    string `json:"from"`
	Message string `json:"body"`
	Type    string `json:"type"`
	Time    int    `json:"time"`
}

type ContactMessage struct {
	Contacts []string `json:"contacts"`
	Type     string   `json:"type"`
}

//middleware
func mongoDB(adress, db string) martini.Handler {
	session, err := mgo.Dial(adress)
	if err != nil {
		panic(err)
	}
	return func(c martini.Context) {
		s := session.Clone()
		c.Map(s.DB(db))
		defer s.Close()
		c.Next()
	}
}
