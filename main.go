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

package main

import (
	"flag"
	"github.com/christopherL91/GoWebSocket/SocketServer"
	"html/template"
	"net/http"
	"runtime"
	"strconv"
)

var (
	port int = 4000
)

func pictures(rw http.ResponseWriter, req *http.Request) {
	req.Header.Add("Content-Type", "image/jpeg")

}

func serveMain(rw http.ResponseWriter, req *http.Request) {
	var template_file, _ = template.ParseFiles("client/src/views/index.html")
	req.Header.Add("Content-Type", "application/javascript")
	template_file.Execute(rw, nil)
}

func init() {
	cores := flag.Int("cores", 1, "The number of cores used")
	flag.Parse()
	runtime.GOMAXPROCS(*cores)
}

func redir(w http.ResponseWriter, req *http.Request) {
	host := req.Host
	http.Redirect(w, req, "https://"+host+":4000/", http.StatusMovedPermanently)
}

func main() {
	server := SocketServer.NewServer("/echo")
	go server.Listen()
	http.HandleFunc("/profilepic/", pictures)
	http.HandleFunc("/", serveMain)
	http.HandleFunc("/public/", func(w http.ResponseWriter, r *http.Request) {
		path := "client/" + r.URL.Path
		http.ServeFile(w, r, path)
	})
	go func() {
		err := http.ListenAndServeTLS(":"+strconv.Itoa(port), "ssl/server.crt", "ssl/server.key", nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()
	if err := http.ListenAndServe(":1337", http.HandlerFunc(redir)); err != nil {
		panic("ListenAndServe error: " + err.Error())
	}
}
