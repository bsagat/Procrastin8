package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Task struct {
	Id         bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Title      string        `bson:"title" json:"title"`
	ActiveDate string        `bson:"activeAt" json:"activeAt"`
	Status     string        `bson:"status" json:"status"`
}
