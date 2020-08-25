package betypes

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type Cashbox struct {
	ID bson.ObjectId			`bson:"_id,omitempty"`
	Type string					`bson:"type"`
	Money float64				`bson:"money"`
	Transactions []Transaction	`bson:"transactions"`
}

type Transaction struct {
	ID bson.ObjectId	`bson:"_id,omitempty"`
	Author string		`bson:"author"`
	Diff float64		`bson:"diff"`
	DataTime time.Time	`bson:"data_time"`
	Comment string		`bson:"comment"`
}