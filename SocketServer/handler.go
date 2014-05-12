package SocketServer

import (
	"code.google.com/p/gcfg"
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

//Struct to hold all the configurations.
type Config struct {
	Server struct {
		Port string
		Bind string
	}
	DB struct {
		Address string
		Name    string
	}
}

func StartServer() {

	//---------------------------------------------------------------
	//Config area
	config := new(Config)
	var port string
	var bindAddress string
	var dbaddress string
	var dbname string

	err := gcfg.ReadFileInto(config, "SocketServer/config.gcfg")

	if err != nil {
		log.Fatal("Could not read and parse config file")
	}

	if len(config.Server.Port) == 0 {
		port = "4000" //Default value if no config file is present.
	} else {
		port = config.Server.Port
	}

	if len(config.Server.Bind) == 0 {
		bindAddress = "127.0.0.1" //Default value if no config file is present.
	} else {
		bindAddress = config.Server.Bind
	}

	if len(config.DB.Address) == 0 {
		dbaddress = "localhost" //Default value if no config file is present.
	} else {
		dbaddress = config.DB.Address
	}

	if len(config.DB.Name) == 0 {
		dbname = "Dogma" //Default value if no config file is present.
	} else {
		dbname = config.DB.Name
	}
	//---------------------------------------------------------------

	m := martini.Classic()
	log.Println("New server started on port " + port)
	//create new cookie store
	store := sessions.NewCookieStore([]byte("mySuperSecretPassword1234"))
	m.Use(martini.Static("bower_components", martini.StaticOptions{
		Prefix: "bower",
	}))
	//Specify new session
	m.Use(sessions.Sessions("Dogma-session", store))

	//specify which database to use.
	m.Use(mongoDB(dbaddress, dbname))

	//Used to print out JSON messages.
	m.Use(render.Renderer())

	//specify session options.
	store.Options(sessions.Options{
		Secure:   true,
		HttpOnly: true,
		MaxAge:   86400, //1 day
	})

	// START OMIT
	m.Post("/login", binding.Form(User{}), func(user User, db *mgo.Database, r render.Render, s sessions.Session) {
		if len(user.Username) > 0 && len(user.Password) > 0 {
			UserID, err := authenticateUser(user.Username, user.Password, db)
			if err != nil {
				r.JSON(401, map[string]interface{}{"error": err.Error()})
				return
			}
			s.Set("UserID", UserID)          //Set session value
			s.Set("Username", user.Username) //Set session value
			r.JSON(200, map[string]interface{}{"_id": UserID, "name": user.Username})
		} else {
			r.JSON(401, map[string]interface{}{"error": "Unauthorized"})
		}
	})
	// END OMIT

	m.Post("/logout", func(s sessions.Session, r render.Render) string {
		s.Delete("UserID")   //Delete session value
		s.Delete("Username") //Delete session value
		r.JSON(200, map[string]interface{}{"status": "logged out"})
		return ""
	})

	m.Get("/status", func(s sessions.Session, r render.Render) {
		userid := s.Get("UserID")
		if userid, ok := userid.(string); ok {
			r.JSON(200, map[string]interface{}{
				"status":   "logged in",
				"loggedIn": true,
			})
		} else {
			r.JSON(401, map[string]interface{}{
				"status":   "logged out",
				"loggedIn": false,
			})
		}
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
			s.Set("UserID", UserID)
			s.Set("Username", user.Username)
			r.JSON(200, map[string]interface{}{"status": "user added"})
		}
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

	m.NotFound(func(rw http.ResponseWriter) {
		rw.Header().Set("Content-Type", "image/jpeg")
		data, _ := os.Open("public/pictures/gophers.jpeg")
		defer data.Close()
		io.Copy(rw, data)
	})

	//---------------------------------------------------------------
	// Secured routes

	server := NewServer(dbaddress, dbname)
	go server.Listen()
	//Only for websocket connection.
	m.Get("/chat", RequireLogin, func(res http.ResponseWriter, req *http.Request, s sessions.Session, db *mgo.Database) {
		username, userid := s.Get("Username"), s.Get("UserID")
		if username, ok := username.(string); ok {
			if userid, ok := userid.(string); ok {
				handler := server.onConnectHandler(username, userid, db)
				handler.ServeHTTP(res, req)
			}
		}
	})

	//---------------------------------------------------------------

	bindAddress = ":" + port
	log.Fatal(http.ListenAndServeTLS(bindAddress,
		"Socketserver/ssl/server.crt",
		"Socketserver/ssl/server.key",
		m))
}
