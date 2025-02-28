package models

import "time"

type Task struct {
	Id         int       `json:"id"`
	Title      string    `json:"title"`
	ActiveDate time.Time `json:"activeAt"`
}
