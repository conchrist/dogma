package Login

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/codegangsta/martini-contrib/sessions"
	"labix.org/v2/mgo"
)

const (
	forms = "/Users/christopher/Documents/Programmering/Go/libs/src/github.com/christopherL91/GoWebSocket/Login/templates"
)

func StartServer() {

	//start martini
	m := martini.Classic()

	store := sessions.NewCookieStore([]byte("mySuperSecretPassword1234"))

	m.Use(sessions.Sessions("login-session", store))

	options := new(sessions.Options)
	options.HttpOnly = true
	options.MaxAge = 86400 * 7
	store.Options(*options)

	/*
	*	use mongoDB as database.
	*	adress,DB
	 */
	m.Use(mongoDB("localhost", "wishDB"))

	//used to display html
	m.Use(render.Renderer(render.Options{
		Directory: forms,
	}))

	m.Get("/", func(r render.Render) {
		r.Redirect("/login", 200)
	})

	m.Get("/login", func(r render.Render, db *mgo.Database, s sessions.Session) string {
		//check if user has a session.
		v := s.Get("userId")
		if v != nil {
			return "already logged in"
		}
		r.HTML(200, "login", nil)
		return ""
	})

	m.Post("/login", binding.Form(User{}), func(userform User, r render.Render, db *mgo.Database, s sessions.Session) (int, string) {
		hashedPass := hashPass(userform.Password)
		Id, err := checkUser(userform.Username, hashedPass, db)
		if err != nil {
			return 401, err.Error()
		}
		s.Set("userId", Id)
		name := userform.Username
		return 200, "logged in as " + name
	})

	m.Get("/userID", func() string {
		return "hej"
	})

	m.Get("/logout", func(s sessions.Session) string {
		s.Delete("userId")
		return "logged out"
	})

	m.Post("/addUser", binding.Form(User{}), func(user User, db *mgo.Database, r render.Render) (int, string) {
		hashedPass := hashPass(user.Password)
		err := addUser(user.Username, hashedPass, db)
		if err != nil {
			return 401, err.Error()
		}
		return 200, "Added user"
	})

	m.Get("/wishes", RequireLogin, func(r render.Render, db *mgo.Database) {
		r.HTML(200, "main", GetAll(db))
	})

	m.Post("/wishes", RequireLogin, binding.Form(Wish{}), func(wish Wish, r render.Render, db *mgo.Database) {
		if len(wish.Name) > 0 && len(wish.Description) > 0 {
			//choose collection.
			db.C("wishes").Insert(wish)
		}
		//load template again.
		r.HTML(200, "main", GetAll(db))
	})

	m.NotFound(func() string {
		return "Sorry, did not find what you were looking for :("
	})

	m.Run()
}
