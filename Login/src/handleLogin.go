package Login

import (
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
	"labix.org/v2/mgo"
	"net/http"
)

const (
	forms = "/Users/christopher/Documents/Programmering/Go/libs/src/github.com/christopherL91/GoWebSocket/Login/templates"
)

func StartServer() {

	store := sessions.NewCookieStore([]byte("1234"))

	//start martini
	m := martini.Classic()

	store.Options(sessions.Options{
		MaxAge: 0,
	})

	m.Use(sessions.Sessions("login_session", store))
	m.Use(sessionauth.SessionUser(GenerateAnonymousUser))
	sessionauth.RedirectUrl = "/login"
	sessionauth.RedirectParam = "new-next"

	/*
	*	use mongoDB as database.
	*	adress,DB
	 */
	m.Use(mongoDB("localhost", "whishDB"))

	//used to display html
	m.Use(render.Renderer(render.Options{
		Directory: forms,
	}))

	m.Get("/", func(r render.Render) {
		r.Redirect("/login", 301)
	})

	m.Get("/login", func(r render.Render, db *mgo.Database) {
		r.HTML(200, "login", nil)
	})

	m.Post("/login", binding.Form(User{}), func(session sessions.Session, userform User, r render.Render, db *mgo.Database, req *http.Request) (int, string) {
		user := User{}
		hashedPass := hashPass(userform.Password)
		err := checkUser(userform.Username, hashedPass, db)
		if err != nil {
			r.Redirect(sessionauth.RedirectUrl)
			return 501, "something wrong happened"
		} else {
			err = sessionauth.AuthenticateSession(session, &user)
			if err != nil {
				r.JSON(500, "wrong")
			}
			params := req.URL.Query()
			redirect := params.Get(sessionauth.RedirectParam)
			r.Redirect(redirect)
			return 200, "logged in"
		}
	})

	/*
	*	function to add user.
	*	example:
	*	curl --header "API-KEY:secretKey" --data "username=christopher&password=kalle" 	http://127.0.0.1:3000/addUser
	*
	 */
	m.Post("/addUser", binding.Form(User{}), func(user User, db *mgo.Database) (int, string) {
		fmt.Println("hello there")
		hashedPass := hashPass(user.Password)
		fmt.Println("hello")
		tmp := addUser(user.Username, hashedPass, db)
		if tmp != nil {
			return 500, tmp.Error()
		}
		return 200, "Added user"
	})

	m.Get("/whishes", sessionauth.LoginRequired, func(r render.Render, db *mgo.Database) {
		r.HTML(200, "main", GetAll(db))
	})

	m.Post("/whishes", sessionauth.LoginRequired, binding.Form(Wish{}), func(wish Wish, r render.Render, db *mgo.Database) {
		if len(wish.Name) > 0 && len(wish.Description) > 0 {
			//choose collection.
			db.C("wishes").Insert(wish)
		}
		//load template again.
		r.HTML(200, "main", GetAll(db))
	})

	m.Get("/logout", sessionauth.LoginRequired, func(session sessions.Session, user sessionauth.User, r render.Render) {
		sessionauth.Logout(session, user)
		r.Redirect("/")
	})

	m.NotFound(func() string {
		return "Sorry, did not find what you were looking for :("
	})

	m.Run()
}
