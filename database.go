package main

import (
	"github.com/globalsign/mgo"
)

var productsCollection *mgo.Collection

func initMongoDB() error {
	session, err := mgo.Dial("mongodb://" + MongoUsername + ":" + MongoPassword + "@" + MongoHostname + ":" + MongoPort)
	if err != nil {
		return err
	}

	productsCollection = session.DB("crm_bot_db").C("products")

	if err = session.Ping(); err != nil {
		return err
	}
	return nil
}
