package database

import (
	"log"

	"../betypes"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var (
	ProductsCollection *mgo.Collection
	UsersCollection    *mgo.Collection
	ProductTypesCollection *mgo.Collection
	CashboxCollection *mgo.Collection
	DailyCashCollection *mgo.Collection
)
type m bson.M

func init() {
	session, err := mgo.Dial("mongodb://" + betypes.MongoUsername + ":" + betypes.MongoPassword + "@" + betypes.MongoHostname + ":" + betypes.MongoPort)
	if err != nil {
		log.Fatal(err)
	}

	ProductsCollection = session.DB("crm_bot_db").C("products")
	UsersCollection = session.DB("crm_bot_db").C("users")
	ProductTypesCollection = session.DB("crm_bot_db").C("prod_types")
	CashboxCollection = session.DB("crm_bot_db").C("cashbox")
	DailyCashCollection = session.DB("crm_bot_db").C("daily_cash")

	if err = session.Ping(); err != nil {
		log.Fatal(err)
	}
}

func MakeTransaction(t *betypes.Transaction) error {
	who := m{"type": "general"}
	pushToArray := m{
		"$push": m{
			"transactions": t},
		"$inc": m{
			"money": t.Diff,
		},	
	}

	err := CashboxCollection.Update(who, pushToArray)
	if err != nil {
		return err
	}
	return nil
}