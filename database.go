package main

import (
	"github.com/globalsign/mgo"
)

var (
	ProductsCollection  *mgo.Collection
	UsersCollection *mgo.Collection
)

func initMongoDB() error {
	session, err := mgo.Dial("mongodb://" + MongoUsername + ":" + MongoPassword + "@" + MongoHostname + ":" + MongoPort)
	if err != nil {
		return err
	}

	ProductsCollection = session.DB("crm_bot_db").C("products")
	UsersCollection = session.DB("crm_bot_db").C("users")

	if err = session.Ping(); err != nil {
		return err
	}
	return nil
}
