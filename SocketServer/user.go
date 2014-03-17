package SocketServer

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"github.com/codegangsta/martini-contrib/sessions"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
)

type User struct {
	UserID   bson.ObjectId `json:"id" bson:"_id"`
	Username string        `form:"username"`
	Password string        `form:"password"`
}

func hashPass(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func validatePass(pass, hash string) bool {
	bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		return false
	}
	return true
}

func authUser(username, password, choice string, db *mgo.Database) (string, error) {
	collection := db.C("Users")

	switch choice {
	case "add":
		user := new(User)
		user.UserID = bson.NewObjectId()
		user.Username = username
		user.Password = password
		if err := collection.Insert(user); err != nil {
			return "", errors.New("could not insert user:" + err.Error())
		}
		return user.UserID.Hex(), nil
	case "check":
		user := new(User)
		err := collection.Find(bson.M{"username": username}).One(user)
		if err != nil {
			//could not find user
			return "", errors.New("Could not find user")
		}
		//the password did not match
		if !validatePass(password, user.Password) {
			return "", errors.New("Could not find user")
		}
		return user.UserID.Hex(), nil
	}
	return "", errors.New("Unhandled case")
}

//rw http.ResponseWriter, req *http.Request, db *mgo.Database, s sessions.Session, r render.Render
//middleware!

func RequireLogin(s sessions.Session, rw http.ResponseWriter,
	req *http.Request) {
	if s.Get("userId") == nil {
		http.Redirect(rw, req, "/login", http.StatusTemporaryRedirect)
	}
}

// func RequireLogin() martini.Handler {
// 	return func(s sessions.Session, context martini.Context, rw http.ResponseWriter, req *http.Request) {
// 		log.Println("Hello there")
// 		id := s.Get("userId")
// 		log.Println(id)

// 		if id == nil {
// 			//http.Redirect(rw, req, "/login", http.StatusFound)
// 			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
// 		}
// 	}
// }
