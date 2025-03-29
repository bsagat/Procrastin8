package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Task struct {
	Id             bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Title          string        `bson:"title" json:"title"`
	ActiveDateStr  string        `json:"activeAt" bson:"activeAtStr"`
	ActiveDateTime time.Time     `bson:"activeAt" json:"-"`
	Status         string        `bson:"status" json:"status"`
}
