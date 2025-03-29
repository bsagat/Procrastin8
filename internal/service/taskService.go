package service

import (
	"TodoApp/internal/domain"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TaskService struct {
	Repo domain.TaskRepoImp
}

func DefaultTaskService(Repo domain.TaskRepoImp) *TaskService {
	return &TaskService{Repo: Repo}
}

var _ domain.TaskServiceImp = (*TaskService)(nil)

// Бизнес логика добавления задачи в базу данных
func (serv *TaskService) CreateTask(task *domain.Task) (int, error) {
	var err error
	task.Status = "active"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	unique, err := serv.Repo.IsTaskUnique(ctx, *task)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !unique {
		return http.StatusBadRequest, errors.New("task data must be unique")
	}
	task.ActiveDateTime, err = time.Parse("2006-01-02", task.ActiveDateStr)
	if err != nil {
		slog.Error("Time Parse error: " + err.Error())
		return http.StatusBadRequest, errors.New("active date field is invalid: time parsing error")
	}
	err = serv.Repo.CreateTask(ctx, task)
	if err != nil {
		return http.StatusBadRequest, err
	}
	return http.StatusCreated, nil
}

// Логика обработки запроса на получение задачи из базы данных
func (serv *TaskService) GetTask(id string) (domain.Task, int, error) {
	var task domain.Task
	objid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return task, http.StatusBadRequest, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	task, err = serv.Repo.GetTask(ctx, objid)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return task, http.StatusNotFound, errors.New("task is not found")
		}
		return task, http.StatusInternalServerError, err
	}
	return task, http.StatusOK, nil
}

// Логика обработки запроса на получение всех задач из базы данных
func (serv *TaskService) GetTasks(status string) ([]domain.Task, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tasks, err := serv.Repo.GetTasks(ctx, status)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if tasks == nil {
		return nil, http.StatusNotFound, nil
	}
	return tasks, http.StatusOK, nil
}

// Логика обработки запроса на обновление задачи из базы данных
func (serv *TaskService) UpdateTask(task domain.Task, id string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return http.StatusBadRequest, err
	}
	task.ActiveDateTime, err = time.Parse("2006-01-02", task.ActiveDateStr)
	if err != nil {
		slog.Error("Time Parse error: " + err.Error())
		return http.StatusBadRequest, errors.New("active date field is invalid: time parsing error")
	}
	err = serv.Repo.UpdateTask(ctx, objID, task)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("task is not found")
	} else if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusNoContent, nil
}

// Логика обработки запроса удаления задачи
func (serv *TaskService) DeleteTask(id string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = serv.Repo.DeleteTask(ctx, objID)
	if err == mongo.ErrNoDocuments {
		return http.StatusNotFound, errors.New("task is not found")
	} else if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusNoContent, nil
}

// Логика обработки запроса изменения статуса задачи
func (serv *TaskService) ChangeStatus(id string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return http.StatusBadRequest, err
	}
	err = serv.Repo.ChangeStatus(ctx, objID)
	if err == mongo.ErrNoDocuments {
		return http.StatusBadRequest, errors.New("task status is already changed")
	} else if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusNoContent, nil
}
