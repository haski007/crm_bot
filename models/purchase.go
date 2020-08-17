package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Purchase struct
type Purchase struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Amount    float64       `bson:"amount"`
	ProductID bson.ObjectId `bson:"product_id"`
	SaleDate  time.Time     `bson:"sale_date"`
}
