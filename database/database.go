package database

import (
	"log"

	"../betypes"
	"github.com/globalsign/mgo"
)

var (
	ProductsCollection *mgo.Collection
	UsersCollection    *mgo.Collection
	ProductTypesCollection *mgo.Collection
)

func init() {
	session, err := mgo.Dial("mongodb://" + betypes.MongoUsername + ":" + betypes.MongoPassword + "@" + betypes.MongoHostname + ":" + betypes.MongoPort)
	if err != nil {
		log.Fatal(err)
	}

	ProductsCollection = session.DB("crm_bot_db").C("products")
	UsersCollection = session.DB("crm_bot_db").C("users")
	ProductTypesCollection = session.DB("crm_bot_db").C("prod_types")

	if err = session.Ping(); err != nil {
		log.Fatal(err)
	}
}
