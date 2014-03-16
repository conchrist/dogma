package SocketServer

import (
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/codegangsta/martini-contrib/sessions"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
)

type User struct {
	UserID   bson.ObjectId `form:"userid"`
	Username string        `form:"username"`
	Password string        `form:"password"`
}

func hashPass(pass string) string {
	hasher := sha512.New()
	hasher.Write([]byte(pass))
	hashedPass := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return hashedPass
}

func authUser(username, password, choice string, db *mgo.Database) (string, error) {
	user := new(User)
	collection := db.C("Users")

	switch choice {
	case "add":
		user.Username = username
		user.Password = password
		user.UserID = bson.NewObjectId()
		if err := collection.Insert(user); err != nil {
			return "", errors.New("could not insert user")
		}
	case "check":
		err := collection.Find(bson.M{"username": username, "password": password}).One(user)
		if err != nil {
			//could not find user
			return "", errors.New("Could not find user")
		}
		//the password did not match
		if user.Password != password {
			return "", errors.New("Could not find user")
		}
	}
	return user.UserID.Hex(), nil
}

//middleware!
func RequireLogin(rw http.ResponseWriter, req *http.Request, db *mgo.Database, s sessions.Session, r render.Render) {
	id := s.Get("userId")

	if id == nil {
		http.Redirect(rw, req, "/login", http.StatusFound)
		return
	} else {
		idString := id.(string)
		user := new(User)
		collection := db.C("Users")
		err := collection.Find(bson.M{"userid": bson.ObjectIdHex(idString)}).One(user)
		if err != nil {
			//the id was not found in the database! We have a scammer here! :P
			http.Redirect(rw, req, "/login", http.StatusFound)
			return
		}
	}
}
