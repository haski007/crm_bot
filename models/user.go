package models

import "github.com/globalsign/mgo/bson"

type User struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	FirstName string        `bson:"first_name"`
	LastName  string        `bson:"last_name"`
	UserName  string        `bson:"username"`
	UserID    int           `bson:"user_id"`
	Status    string        `bson:"status"`
}
