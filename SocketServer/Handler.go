package SocketServer

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/codegangsta/martini-contrib/sessions"
	"github.com/russross/blackfriday"
	"io"
	"io/ioutil"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"os"
)

type Error struct {
	ErrorData string
}

var (
	errMessage = new(Error)
)

func init() {
	server := NewServer("/chatroom")
	go server.Listen()
}

func StartServer() {

	errMessage.ErrorData = "Fill in missing data"

	m := martini.Classic()
	store := sessions.NewCookieStore([]byte("mySuperSecretPassword1234"))

	m.Use(sessions.Sessions("login-session", store))
	m.Use(mongoDB("localhost", "Golang"))
	//used to display html
	m.Use(render.Renderer(render.Options{
		Directory: "public/templates",
	}))

	options := new(sessions.Options)
	options.Secure = true
	options.HttpOnly = true
	options.MaxAge = 86400 //1 day
	store.Options(*options)

	m.Get("/", func(rw http.ResponseWriter, req *http.Request) {
		http.Redirect(rw, req, "/login", http.StatusMovedPermanently)
	})

	m.Get("/login", func(r render.Render, s sessions.Session) string {
		v := s.Get("userId")
		if v != nil {
			return "already logged in"
		}
		r.HTML(200, "login", nil)
		return ""
	})

	m.Post("/login", binding.Form(User{}), func(user User, db *mgo.Database, r render.Render, s sessions.Session) (int, string) {
		if len(user.Username) > 0 && len(user.Password) > 0 {
			hashedPass := hashPass(user.Password)
			ID, err := authUser(user.Username, hashedPass, "check", db)
			if err != nil {
				return 401, err.Error()
			}
			s.Set("userId", ID)
		} else {
			r.HTML(200, "login", errMessage)
			return 401, ""
		}
		return 200, "logged in"
	})

	m.Get("/logout", func(s sessions.Session) string {
		s.Delete("userId")
		return "logged out"
	})

	m.Get("/chat", RequireLogin, func(r render.Render, req *http.Request) {
		r.HTML(200, "chat", nil)
	})

	m.Get("/newuser", func(r render.Render) {
		r.HTML(200, "new", nil)
	})

	m.Post("/newuser", binding.Form(User{}), func(user User, r render.Render, db *mgo.Database, s sessions.Session) (int, string) {
		if len(user.Username) > 0 && len(user.Password) > 0 {
			hashedPass := hashPass(user.Password)
			ID, err := authUser(user.Username, hashedPass, "add", db)
			if err != nil {
				return 401, err.Error()
			}
			s.Set("userId", ID)
		} else {
			r.HTML(200, "new", errMessage)
			return 401, ""
		}
		return 200, "Successfully addded new user"
	})

	m.NotFound(func(rw http.ResponseWriter) {
		rw.Header().Set("Content-Type", "image/jpeg")
		data, _ := os.Open("public/pictures/gopher.jpeg")
		defer data.Close()
		io.Copy(rw, data)
	})

	m.Get("/license", func(rw http.ResponseWriter) {
		input, err := ioutil.ReadFile("public/markdown/license.md")
		if err != nil {
			log.Println(err.Error())
		}
		markdown := blackfriday.MarkdownCommon(input)
		rw.Write(markdown)
	})

	m.Get("/about", func(rw http.ResponseWriter) {
		input, err := ioutil.ReadFile("public/markdown/about.md")
		if err != nil {
			log.Println(err.Error())
		}
		markdown := blackfriday.MarkdownCommon(input)
		rw.Write(markdown)
	})

	log.Fatal(http.ListenAndServeTLS(":4000", "SocketServer/ssl/server.crt", "SocketServer/ssl/server.key", m))
}
