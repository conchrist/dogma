package Login

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"labix.org/v2/mgo"
	"net/http"
)

const (
	forms = "/Users/christopher/Documents/Programmering/Go/libs/src/github.com/christopherL91/GoWebSocket/Login/templates"
)

func StartServer() {
	//start martini
	m := martini.Classic()

	//used to display html
	m.Use(render.Renderer(render.Options{
		Directory: forms,
	}))

	m.Use(martini.Static("pictures"))
	/*
	*	use mongoDB as database.
	*	adress,DB
	 */
	m.Use(mongoDB("localhost", "whishDB"))

	m.Get("/login", func(r render.Render, db *mgo.Database) {
		r.HTML(200, "login", nil)
	})

	m.Post("/login", binding.Form(User{}), func(user User, r render.Render, db *mgo.Database) (int, string) {
		hashedPass := hashPass(user.Password)
		err := checkUser(user.Username, hashedPass, db)
		if err != nil {
			return 401, err.Error()
		}
		//add session
		return 200, "logged in"
	})

	/*
	*	function to add user.
	*	example:
	*	curl --header "API-KEY:secretKey" --data "username=christopher&password=kalle" 	http://127.0.0.1:3000/addUser
	*
	 */
	m.Post("/addUser", Auth, binding.Form(User{}), func(user User, db *mgo.Database) (int, string) {
		hashedPass := hashPass(user.Password)
		tmp := addUser(user.Username, hashedPass, db)
		if tmp != nil {
			return 500, tmp.Error()
		}
		return 200, "Added user"
	})

	m.Get("/whishes", func(r render.Render, db *mgo.Database) {
		r.HTML(200, "main", GetAll(db))
	})

	m.Post("/whishes", binding.Form(Wish{}), func(wish Wish, r render.Render, db *mgo.Database) {
		if len(wish.Name) > 0 && len(wish.Description) > 0 {
			//choose collection.
			db.C("Whishes").Insert(wish)
			//load template again.
		}
		r.HTML(200, "main", GetAll(db))
	})

	m.NotFound(func() string {
		return "Sorry, did not find what you were looking for :("
	})

	m.Run()
}

func Auth(rw http.ResponseWriter, req *http.Request) {
	if req.Header.Get("API-KEY") != "secretKey" {
		http.Error(rw, "Provide right API-KEY", 401)
	}
}
