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

func StartServer() {

	m := martini.Classic()
	store := sessions.NewCookieStore([]byte("mySuperSecretPassword1234"))
	m.Use(martini.Static("bower_components", martini.StaticOptions{
		Prefix: "bower",
	}))

	m.Use(sessions.Sessions("login-session", store))
	m.Use(mongoDB("localhost", "Golang"))
	//used to render html or json
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

	m.Get("/login", func(r render.Render, s sessions.Session) {
		UserID := s.Get("UserID")
		Username := s.Get("Username")
		if UserID != nil {
			r.JSON(200, map[string]interface{}{"_id": UserID, "name": Username})
			return
		}
		r.HTML(200, "login", nil)
		return
	})

	m.Post("/login", binding.Form(User{}), func(user User, db *mgo.Database, r render.Render, s sessions.Session) {
		if len(user.Username) > 0 && len(user.Password) > 0 {
			UserID, err := authenticateUser(user.Username, user.Password, db)
			if err != nil {
				return 401, err.Error()
			}
			s.Set("UserID", UserID)
			s.Set("Username", user.Username)
			r.JSON(200, map[string]interface{}{"_id": UserID, "name": user.Username})
			return
		} else {
			r.JSON(401, map[string]interface{}{"error": "Unauthorized"})
			return
		}
	})

	m.Get("/logout", func(s sessions.Session, r render.Render) string {
		s.Delete("UserID")
		r.JSON(200, map[string]interface{}{"status": "logged out"})
		return ""
	})

	m.Get("/users/new", func(r render.Render) {
		r.HTML(200, "new", nil)
	})

	m.Post("/users", binding.Form(User{}), func(user User, r render.Render, db *mgo.Database, s sessions.Session) {
		if len(user.Username) > 0 && len(user.Password) > 0 {
			passwordHash, err := hashPass(user.Password)
			if err != nil {
				r.JSON(401, map[string]interface{}{"error": err.Error()})
				return
			}
			UserID, err := addUser(user.Username, passwordHash, db)
			if err != nil {
				r.JSON(401, map[string]interface{}{"error": err.Error()})
				return
			}
			s.Set("userId", UserID)
			r.JSON(200, map[string]interface{}{"status": "user added"})
			return 200, ""
		}
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
	//---------------------------------------------------------------
	// Secured routes

	server := NewServer()
	go server.Listen()
	m.Get("/chat", RequireLogin, func(res http.ResponseWriter, req *http.Request) {
		handler := server.onConnectHandler()
		handler.ServeHTTP(res, req)
	})

	//---------------------------------------------------------------

	log.Fatal(http.ListenAndServeTLS(":4000", "SocketServer/ssl/server.crt", "SocketServer/ssl/server.key", m))
}
