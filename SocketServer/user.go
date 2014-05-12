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

func addUser(username, hashedPassword string, db *mgo.Database) (string, error) {
	collection := db.C("Users")
	user := new(User)
	user.UserID = bson.NewObjectId()
	user.Username = username
	user.Password = hashedPassword
	if err := collection.Insert(user); err != nil {
		return "", errors.New("could not insert user:" + err.Error())
	}
	return user.UserID.Hex(), nil
}

func authenticateUser(username, password string, db *mgo.Database) (string, error) {
	collection := db.C("Users")
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

func RequireLogin(s sessions.Session, rw http.ResponseWriter,
	req *http.Request) {
	if s.Get("UserID") == nil {
		http.Redirect(rw, req, "/", http.StatusTemporaryRedirect)
	}
}
