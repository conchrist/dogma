package Login

import (
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"github.com/martini-contrib/sessionauth"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type User struct {
	Id            string `form:"_id"`
	Username      string `form:"username"`
	Password      string `form:"password"`
	authenticated bool   `form:"authenticated"`
}

func hashPass(pass string) string {
	hasher := sha512.New()
	hasher.Write([]byte(pass))
	hashedPass := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return hashedPass
}

/*
* db.Users.find({username:"Kalle", password:1234})
 */
func checkUser(username, password string, db *mgo.Database) error {
	user := new(User)
	coll := db.C("Users")
	err := coll.Find(bson.M{"username": username, "password": password}).One(user)
	if err != nil {
		//could not find user
		return errors.New("Could not find user")
	}
	if user.Password != password {
		return errors.New("Could not find user")
	}
	return nil
}

func addUser(username, password string, db *mgo.Database) error {

	user := new(User)
	user.Username = username
	user.Password = password
	user.authenticated = true

	//set collection in database.
	collection := db.C("Users")

	if err := collection.Insert(user); err != nil {
		return errors.New("could not insert user")
	}
	return nil
}

//Part of the sessionauth API
/*------------------------------------
 */

func (u *User) IsAuthenticated() bool {
	return u.authenticated
}

func (u *User) Login() {
	u.authenticated = true
}

func (u *User) Logout() {
	u.authenticated = false
}

func (u *User) UniqueId() interface{} {
	return u.Id
}

func (u *User) GetById(id interface{}) error {
	return nil
}

func GenerateAnonymousUser() sessionauth.User {
	return &User{}
}

/*
----------------------------------------
*/
