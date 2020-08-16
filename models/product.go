package models

import (
	"github.com/globalsign/mgo/bson"
)

// Product struct.
type Product struct {
	ID    bson.ObjectId `json:"-" bson:"_id"`
	Name  string             `json:"name" bson:"name"`
	Type  string             `json:"type" bson:"type"`
	Price float64            `json:"price" bson:"price"`
}
