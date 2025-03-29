package handlers

import (
	"TodoApp/internal/domain"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Taskhandler struct {
	Serv domain.TaskServiceImp
}

func DefaultTaskHandler(serv domain.TaskServiceImp) *Taskhandler {
	return &Taskhandler{Serv: serv}
}

var _ domain.TaskHandlerImp = (*Taskhandler)(nil)

// Обработчик логики создания задач
func (h *Taskhandler) NewTaskHandler(c *gin.Context) {
	// Парсинг JSON запроса
	var task domain.Task
	err := c.BindJSON(&task)
	if err != nil {
		slog.Error("JSON Bind error: ",
			slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Валидация задачи
	err = h.ValidateTask(task)
	if err != nil {
		slog.Error("Task validation error: ",
			slog.String("error", err.Error()),
			slog.String("id", task.Id.String()),
			slog.String("active date", task.ActiveDateStr))
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	// Вызов бизнес логики создания задачи
	code, err := h.Serv.CreateTask(&task)
	if err != nil {
		slog.Error("Task Create error: ",
			slog.String("error", err.Error()),
			slog.String("id", task.Id.String()),
			slog.String("active date", task.ActiveDateStr))
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}
	slog.Info("Task created succesfully:", slog.String("Task id", task.Id.String()))
	c.JSON(code, gin.H{"inserted_id": task.Id})
}

// Обработчик для получения задачи по уникальному ID
func (h *Taskhandler) GetTaskHandler(c *gin.Context) {
	// Чтение параметров запроса
	id := c.Param("id")
	if id == "" {
		slog.Error("Parameter read error: ", slog.String("error", "id parameter is empty"))
		c.JSON(http.StatusBadRequest, gin.H{"message": "id field is empty"})
		return
	}
	// Вызов сервисной функции
	task, code, err := h.Serv.GetTask(id)
	if err != nil {
		slog.Error("Get task list error: ", slog.String("error: ", err.Error()))
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}

	slog.Info("Task received succesfully:", slog.String("Task id", id))
	c.JSON(code, task)
}

// Обработчик для получения информации о всех задачах
func (h *Taskhandler) TaskListsHandler(c *gin.Context) {
	// Чтение параметров запроса
	status := c.DefaultQuery("status", "active")
	if status != "active" && status != "done" {
		slog.Error("Query reading error: status parameter is invalid")
		c.JSON(http.StatusBadRequest, gin.H{"message": "status parameter is invalid"})
		return
	}
	// Вызов сервисной функции
	tasks, code, err := h.Serv.GetTasks(status)
	if err != nil {
		slog.Error("Get task list error: ", slog.String("error: ", err.Error()))
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}
	slog.Info("Task lists handled succesfully")
	c.JSON(code, tasks)
}

// Обработчик логики обновления задач
func (h *Taskhandler) UpdateTaskHandler(c *gin.Context) {
	// Чтение параметров запроса
	id := c.Param("id")
	if id == "" {
		slog.Error("Task Update error: ", slog.String("error", "id parameter is empty"))
		c.JSON(http.StatusBadRequest, gin.H{"message": "id field is empty"})
		return
	}

	// Парсинг JSON  запросы
	var task domain.Task
	err := c.BindJSON(&task)
	if err != nil {
		slog.Error("JSON Bind error: ",
			slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Валидация задачи
	err = h.ValidateTask(task)
	if err != nil {
		slog.Error("Task validation error: ",
			slog.String("error", err.Error()),
			slog.String("id", task.Id.String()),
			slog.String("active date", task.ActiveDateStr))
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	// Вызов бизнес логики обновления задач
	code, err := h.Serv.UpdateTask(task, id)
	if err != nil {
		slog.Error("Task Update error: ", slog.String("error", err.Error()))
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}
	slog.Info("Task updated succesfully:", slog.String("Task id", id))
	c.JSON(code, gin.H{"updated_id": id})
}

// Обработчик логики удаления задач
func (h *Taskhandler) DeleteTaskHandler(c *gin.Context) {
	// Чтение параметров запроса
	id := c.Param("id")
	if id == "" {
		slog.Error("Parameter read error: ", slog.String("error", "id parameter is empty"))
		c.JSON(http.StatusBadRequest, gin.H{"message": "id field is empty"})
		return
	}
	// Вызов бизнес логики удаления задачи
	code, err := h.Serv.DeleteTask(id)
	if err != nil {
		slog.Error("Task Delete error", slog.String("error", err.Error()))
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}
	slog.Info("Task Deleted succesfully:", slog.String("Task id", id))
	c.JSON(code, gin.H{"deleted_id": id})
}

// Обработчик логики завершения задачи
func (h *Taskhandler) FinishTaskHandler(c *gin.Context) {
	// Чтение параметров запроса
	id := c.Param("id")
	if id == "" {
		slog.Error("Parameter read error: ", slog.String("error", "id parameter is empty"))
		c.JSON(http.StatusBadRequest, gin.H{"message": "id field is empty"})
		return
	}
	// Вызов сервисной функции
	code, err := h.Serv.ChangeStatus(id)
	if err != nil {
		slog.Error("Task Finish error: ", slog.String("error", "id parameter is empty"))
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}
	slog.Info("Task finished succesfully:", slog.String("Task id", id))
	c.JSON(code, gin.H{"message": "task status changed"})
}

// Логика валидации задачи
func (h *Taskhandler) ValidateTask(task domain.Task) error {
	if !task.Id.IsZero() {
		return errors.New("id string should be empty")
	}
	if task.ActiveDateStr == "" {
		return errors.New("active date field is not exist")
	}

	activeDate, err := time.Parse("2006-01-02", task.ActiveDateStr)
	if err != nil {
		slog.Error("Time Parse error: " + err.Error())
		return errors.New("active date field is invalid: time parsing error")
	}

	// Поле ActiveAt валидно если ( 1925-01-01 < ActiveAt < 2100-01-01 )
	if activeDate.Before(time.Date(1925, 1, 1, 0, 0, 0, 0, time.UTC)) || activeDate.After(time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)) {
		return errors.New("active date field is invalid, 1925-01-01 < ActiveAt < 2100-01-01")
	}

	if task.Status != "" {
		return errors.New("status field must be empty")
	}
	if task.Title == "" {
		return errors.New("title field is not exist")
	}
	if len(task.Title) > 200 {
		return errors.New("title lenght is greater than 200")
	}
	return nil
}
