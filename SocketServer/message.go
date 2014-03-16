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
	"github.com/codegangsta/martini"
	"labix.org/v2/mgo"
)

type MessageStruct struct {
	From    string `json:"from"`
	Message string `json:"body"`
	Type    string `json:"type"`
	Time    int    `json:"time"`
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
