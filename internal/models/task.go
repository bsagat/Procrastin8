package models

import "time"

type Task struct {
	Id         int       `bson:"id"`
	Title      string    `bson:"title"`
	ActiveDate time.Time `bson:"activeAt"`
	Status     string    `bson:"status"`
}
