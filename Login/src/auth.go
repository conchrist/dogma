package Login

import (
	"crypto/sha512"
	"encoding/base64"
	//"github.com/codegangsta/martini"
	"errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type User struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func hashPass(pass string) string {
	hasher := sha512.New()
	hasher.Write([]byte(pass))
	hashedPass := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return hashedPass
}

/*
* db.Users.find({username:"Kalle"} && {password:1234})
 */
func checkUser(username, password string, db *mgo.Database) error {
	user := new(User)
	coll := db.C("Users")
	err := coll.Find(bson.M{"username": username}).One(user)
	if err != nil {
		//could not find user
		return errors.New("Could not find user")
	}
	if user.Password != password {
		return errors.New("Wrong password!")
	}
	return nil
}

func addUser(username, password string, db *mgo.Database) error {

	user := new(User)
	user.Username = username
	user.Password = password
	//set collection
	coll := db.C("Users")

	//query the DB to see if the user already exists.
	num, _ := coll.Find(bson.M{"username": user.Username}).Count()
	if num != 0 {
		return errors.New("User already exists")
	}
	if err := coll.Insert(user); err != nil {
		return errors.New("could not insert value")
	}
	//everything was ok.
	return nil
}
