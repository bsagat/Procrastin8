package test

import (
	"errors"
	"net/http"
	"testing"

	"TodoApp/internal/domain"
	"TodoApp/internal/service"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Для быстрого запуска тестов лучше использовать:  go test -parallel 4 ./tests

type TestCase struct {
	Name     string
	Arg      domain.Task
	Code     int
	Expected error
}

func TestCreateTask(t *testing.T) {
	tests := []TestCase{
		{
			Name:     "Валидная задача",
			Arg:      domain.Task{Title: "", ActiveDateStr: "2025-08-03"},
			Code:     http.StatusCreated,
			Expected: nil,
		},
		{
			Name:     "Не уникальная задача",
			Arg:      domain.Task{Title: "Взять кота с приюта", ActiveDateStr: "2025-09-09"},
			Code:     http.StatusBadRequest,
			Expected: domain.ErrNotUniqueTask,
		},
		{
			Name:     "Неправильный Active Date",
			Arg:      domain.Task{Title: "Test active date", ActiveDateStr: "20-20-20"},
			Code:     http.StatusBadRequest,
			Expected: domain.ErrActiveDateError,
		},
	}
	id := bson.NewObjectID()
	repo := &MockRepo{tasks: map[bson.ObjectID]domain.Task{
		id: {
			Title:         "Взять кота с приюта",
			ActiveDateStr: "2025-09-09",
		},
	}}
	serv := service.DefaultTaskService(repo)
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			code, err := serv.CreateTask(&tt.Arg)
			if code != tt.Code {
				t.Errorf("Create Task testing fail: expected: %d,found: %d", tt.Code, code)
			}
			if !errors.Is(tt.Expected, err) {
				t.Errorf("Create Task testing fail: expected: %s,found: %s", tt.Expected.Error(), err.Error())
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	id := bson.NewObjectID()
	tests := []TestCase{
		{
			Name:     "Валидный ID",
			Arg:      domain.Task{Id: id},
			Code:     http.StatusOK,
			Expected: nil,
		},
		{
			Name:     "Не существующий ID",
			Arg:      domain.Task{Id: bson.NewObjectID()},
			Code:     http.StatusNotFound,
			Expected: domain.ErrTaskNotFound,
		},
	}
	repo := &MockRepo{tasks: map[bson.ObjectID]domain.Task{
		id: {
			Title:         "Взять кота с приюта",
			ActiveDateStr: "2025-09-09",
		},
	}}
	serv := service.DefaultTaskService(repo)
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, code, err := serv.GetTask(tt.Arg.Id.Hex())
			if code != tt.Code {
				t.Errorf("Get Task testing fail: expected: %d,found: %d", tt.Code, code)
			}
			if err != tt.Expected {
				t.Errorf("Get Task testing fail: expected: %s,found: %s", tt.Expected.Error(), err.Error())
			}
		})
	}
}

func TestGetTasks(t *testing.T) {
	tests := []TestCase{
		{
			Name:     "Валидный тест кейс",
			Arg:      domain.Task{Status: "active"},
			Code:     http.StatusOK,
			Expected: nil,
		},
		{
			Name:     "Неправильный статус",
			Arg:      domain.Task{Status: "invalid status"},
			Code:     http.StatusBadRequest,
			Expected: domain.ErrInvalidStatus,
		},
	}
	id := bson.NewObjectID()
	repo := &MockRepo{tasks: map[bson.ObjectID]domain.Task{
		id: {
			Title:         "Взять кота с приюта",
			ActiveDateStr: "2023-09-09",
		},
	}}
	serv := service.DefaultTaskService(repo)
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, code, err := serv.GetTasks(tt.Arg.Status)
			if code != tt.Code {
				t.Errorf("Get Tasks testing fail: expected: %d,found: %d", tt.Code, code)
			}
			if err != tt.Expected {
				t.Errorf("Get Tasks testing fail: expected: %s,found: %s", tt.Expected.Error(), err.Error())
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	id := bson.NewObjectID()
	tests := []TestCase{
		{
			Name: "Валидный ID",
			Arg: domain.Task{
				Id:            id,
				Title:         "Взять корм",
				ActiveDateStr: "2025-09-08",
			},
			Code:     http.StatusNoContent,
			Expected: nil,
		},
		{
			Name: "Не существующий ID",
			Arg: domain.Task{
				Id:            bson.NewObjectID(),
				Title:         "Взять корм",
				ActiveDateStr: "2025-09-08",
			},
			Code:     http.StatusNotFound,
			Expected: domain.ErrTaskNotFound,
		},
		{
			Name:     "Неправильный Active Date",
			Arg:      domain.Task{Title: "Test active date", ActiveDateStr: "20-20-20"},
			Code:     http.StatusBadRequest,
			Expected: domain.ErrActiveDateError,
		},
	}
	repo := &MockRepo{tasks: map[bson.ObjectID]domain.Task{
		id: {
			Title:         "Взять кота с приюта",
			ActiveDateStr: "2025-09-09",
		},
	}}
	serv := service.DefaultTaskService(repo)
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			code, err := serv.UpdateTask(tt.Arg, tt.Arg.Id.Hex())
			if code != tt.Code {
				t.Errorf("Update task testing fail: expected: %d,found: %d", tt.Code, code)
			}
			if err != tt.Expected {
				t.Errorf("Update task testing fail: expected: %s,found: %s", tt.Expected.Error(), err.Error())
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	id := bson.NewObjectID()
	tests := []TestCase{
		{
			Name:     "Валидный ID",
			Arg:      domain.Task{Id: id},
			Code:     http.StatusNoContent,
			Expected: nil,
		},
		{
			Name:     "Не существующий ID",
			Arg:      domain.Task{Id: bson.NewObjectID()},
			Code:     http.StatusNotFound,
			Expected: domain.ErrTaskNotFound,
		},
	}
	repo := &MockRepo{tasks: map[bson.ObjectID]domain.Task{
		id: {
			Title:         "Взять кота с приюта",
			ActiveDateStr: "2025-09-09",
		},
	}}
	serv := service.DefaultTaskService(repo)
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			code, err := serv.DeleteTask(tt.Arg.Id.Hex())
			if code != tt.Code {
				t.Errorf("Delete task testing fail: expected: %d,found: %d", tt.Code, code)
			}
			if err != tt.Expected {
				t.Errorf("Delete Task testing fail: expected: %s,found: %s", tt.Expected.Error(), err.Error())
			}
		})
	}
}

func TestChangeStatus(t *testing.T) {
	id := bson.NewObjectID()
	id2 := bson.NewObjectID()
	tests := []TestCase{
		{
			Name: "Валидный task",
			Arg: domain.Task{
				Id: id,
			},
			Code:     http.StatusNoContent,
			Expected: nil,
		},
		{
			Name: "ID которого нету",
			Arg: domain.Task{
				Id: bson.NewObjectID(),
			},
			Code:     http.StatusBadRequest,
			Expected: domain.ErrTaskNotFound,
		},
		{
			Name: "Изменить уже выполненную задачу",
			Arg: domain.Task{
				Id: id2,
			},
			Code:     http.StatusBadRequest,
			Expected: domain.ErrTaskChanged,
		},
	}
	repo := &MockRepo{tasks: map[bson.ObjectID]domain.Task{
		id: {
			Title:         "Взять кота с приюта",
			ActiveDateStr: "2025-09-09",
			Status:        "active",
		},
		id2: {
			Title:         "Взять собакена с приюта",
			ActiveDateStr: "2025-09-09",
			Status:        "done",
		},
	}}
	serv := service.DefaultTaskService(repo)
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			code, err := serv.ChangeStatus(tt.Arg.Id.Hex())
			if code != tt.Code {
				t.Errorf("Change task testing fail: expected: %d,found: %d", tt.Code, code)
			}
			if err != tt.Expected {
				t.Errorf("Change task testing fail: expected: %s,found: %s", tt.Expected.Error(), err.Error())
			}
		})
	}
}
