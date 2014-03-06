package Login

import (
	"github.com/codegangsta/martini"
	"labix.org/v2/mgo"
)

type Wish struct {
	Name        string `form:"name"`
	Description string `form:"description"`
}

//create handler for martini.
func mongoDB(adress, db string) martini.Handler {
	session, err := mgo.Dial(adress)
	if err != nil {
		panic(err)
	}
	return func(c martini.Context) {
		s := session.Clone()
		c.Map(s.DB(db))
		// defer s.Close()
		// c.Next()
	}
}

//return all the wishes from the database.
func GetAll(db *mgo.Database) []Wish {
	var wishlist []Wish
	db.C("wishes").Find(nil).All(&wishlist)
	return wishlist
}
