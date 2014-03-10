package Login

import (
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
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

func checkUser(username, password string, db *mgo.Database) (string, error) {
	user := new(User)
	collection := db.C("users")
	err := collection.Find(bson.M{"username": username, "password": password}).One(user)
	if err != nil {
		//could not find user
		return "", errors.New("Could not find user")
	}
	//the password did not match
	if user.Password != password {
		return "", errors.New("Could not find user")
	}
	return user.UserID.Hex(), nil
}

func addUser(username, password string, db *mgo.Database) error {

	user := new(User)
	user.Username = username
	user.Password = password
	user.UserID = bson.NewObjectId()

	//set collection in database.
	collection := db.C("users")

	if err := collection.Insert(user); err != nil {
		return errors.New("could not insert user")
	}
	return nil
}

//middleware!
func RequireLogin(rw http.ResponseWriter, req *http.Request, db *mgo.Database, s sessions.Session, r render.Render) {
	user := new(User)
	collection := db.C("users")
	id := s.Get("userId")

	if id == nil {
		fmt.Println("no cookie")
		http.Redirect(rw, req, "http://nedry.ytmnd.com/", http.StatusFound)
		return
	} else {
		idString := id.(string)
		err := collection.Find(bson.M{"userid": bson.ObjectIdHex(idString)}).One(user)
		if err != nil {
			//the id was not found in the database! We have a scammer here! :P
			http.Redirect(rw, req, "http://nedry.ytmnd.com/", http.StatusFound)
			return
		}
	}
}
