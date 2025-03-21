package service

import (
	repo "TodoApp/internal/dal"
	"TodoApp/internal/models"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestCreateTask(t *testing.T) {
	db, err := mongo.Connect()
	repo := repo.DefaultTaskRepository(db)
	serv := DefaultTaskService(*repo)
	in := models.Task{
		Id:    bson.NewObjectID(),
		Title: "",
	}
	_, err = serv.CreateTask(&in)
	if err != nil {
		t.Error(err)

	}
}
