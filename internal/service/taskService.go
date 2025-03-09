package service

import (
	repo "TodoApp/internal/dal"
	"TodoApp/internal/models"
	"context"
	"errors"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TaskService struct {
	Repo repo.TaskRepository
}

func DefaultTaskService(Repo repo.TaskRepository) *TaskService {
	return &TaskService{Repo: Repo}
}

func (serv *TaskService) CreateTask(task models.Task) (int, error) {
	task.Status = "active"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	unique, err := serv.Repo.IsTaskUnique(ctx, task)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !unique {
		return http.StatusBadRequest, errors.New("task data must be unique")
	}

	err = serv.Repo.CreateTask(ctx, task)
	if err != nil {
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func (serv *TaskService) GetTasks(status string) ([]models.Task, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tasks, err := serv.Repo.GetTasks(ctx, status)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if tasks == nil {
		return nil, http.StatusNotFound, mongo.ErrNoDocuments
	}
	return tasks, http.StatusOK, nil
}

func (serv *TaskService) ChangeStatus(id string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return http.StatusBadRequest, err
	}
	err = serv.Repo.ChangeStatus(ctx, objID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusNoContent, nil
}

func (serv *TaskService) DeleteTask(id string) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = serv.Repo.DeleteTask(ctx, objID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusNoContent, nil
}

func (serv *TaskService) UpdateTask(task models.Task, id string) (int, error) {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// TODO dodelat update

}
