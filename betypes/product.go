package betypes

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Product struct.
type Product struct {
	ID        bson.ObjectId `json:"-" bson:"_id,omitempty"`
	Name      string        `json:"name" bson:"name"`
	Type      string        `json:"type" bson:"type"`
	Price     float64       `json:"price" bson:"price"`
	PrimeCost float64 		`json:"prime-cost" bson:"prime_cost"`
	InStorage float64 		`json:"in-storage" bson:"in_storage"`
	Purchases []Purchase    `json:"-" bson:"purchases"`
}

// Purchase struct
type Purchase struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Amount   float64       `bson:"amount"`
	SaleDate time.Time     `bson:"sale_date"`
	Seller   string           `bson:"seller"`
}
