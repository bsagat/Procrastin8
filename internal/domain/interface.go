package domain

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TaskHandlerImp interface {
	DeleteTaskHandler(c *gin.Context)
	FinishTaskHandler(c *gin.Context)
	GetTaskHandler(c *gin.Context)
	NewTaskHandler(c *gin.Context)
	TaskListsHandler(c *gin.Context)
	UpdateTaskHandler(c *gin.Context)
	ValidateTask(task Task) error
}

type TaskServiceImp interface {
	CreateTask(task *Task) (int, error)
	DeleteTask(id string) (int, error)
	GetTask(id string) (Task, int, error)
	GetTasks(status string) ([]Task, int, error)
	UpdateTask(task Task, id string) (int, error)
	ChangeStatus(id string) (int, error)
}

type TaskRepoImp interface {
	CreateTask(ctx context.Context, task *Task) error
	DeleteTask(ctx context.Context, id bson.ObjectID) error
	GetTask(ctx context.Context, id bson.ObjectID) (Task, error)
	GetTasks(ctx context.Context, status string) ([]Task, error)
	IsTaskUnique(ctx context.Context, task Task) (bool, error)
	UpdateTask(ctx context.Context, id bson.ObjectID, task Task) error
	ChangeStatus(ctx context.Context, id bson.ObjectID) error
}
