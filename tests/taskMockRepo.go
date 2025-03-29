package test

import (
	"context"

	"TodoApp/internal/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MockRepo struct {
	tasks map[bson.ObjectID]domain.Task
}

func (repo *MockRepo) CreateTask(ctx context.Context, task *domain.Task) error {
	// task.Id = bson.NewObjectID()
	// repo.tasks[task.Id] = *task
	return nil
}

func (repo *MockRepo) GetTask(ctx context.Context, id bson.ObjectID) (domain.Task, error) {
	task, ok := repo.tasks[id]
	if !ok {
		return task, mongo.ErrNoDocuments
	}
	return task, nil
}

func (repo *MockRepo) GetTasks(ctx context.Context, status string) ([]domain.Task, error) {
	var tasks []domain.Task
	for _, task := range repo.tasks {
		tasks = append(tasks, task)
	}
	if tasks == nil {
		return nil, mongo.ErrNoDocuments
	}
	return tasks, nil
}

func (repo *MockRepo) IsTaskUnique(ctx context.Context, task domain.Task) (bool, error) {
	for _, value := range repo.tasks {
		if value.ActiveDateStr == task.ActiveDateStr && value.Title == task.Title {
			return false, nil
		}
	}
	return true, nil
}

func (repo *MockRepo) UpdateTask(ctx context.Context, id bson.ObjectID, task domain.Task) error {
	_, ok := repo.tasks[id]
	if !ok {
		return mongo.ErrNoDocuments
	}
	// repo.tasks[id] = task
	return nil
}

func (repo *MockRepo) ChangeStatus(ctx context.Context, id bson.ObjectID) error {
	task, ok := repo.tasks[id]
	if !ok {
		return mongo.ErrNoDocuments
	}
	if task.Status == "done" {
		return domain.ErrTaskChanged
	}
	return nil
}

func (repo *MockRepo) DeleteTask(ctx context.Context, id bson.ObjectID) error {
	_, ok := repo.tasks[id]
	if !ok {
		return mongo.ErrNoDocuments
	}
	// delete(repo.tasks, id)
	return nil
}
