package handlers

import (
	"TodoApp/internal/models"
	"TodoApp/internal/service"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Taskhandler struct {
	Serv service.TaskService
}

func DefaultTaskHandler(serv service.TaskService) *Taskhandler {
	return &Taskhandler{Serv: serv}
}

// Обработчик логики создания задач
func (h *Taskhandler) NewTaskHandler(c *gin.Context) {
	// Парсинг JSON запроса
	var task models.Task
	err := c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Валидация задачи
	err = h.ValidateTask(task)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	// Вызов бизнес логики создания задачи
	code, err := h.Serv.CreateTask(&task)
	if err != nil {
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}
	c.JSON(code, gin.H{"inserted_id": task.Id})
}

// Обработчик логики обновления задач
func (h *Taskhandler) UpdateTaskHandler(c *gin.Context) {
	// Чтение параметров запроса
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id field is empty"})
		return
	}

	// Парсинг JSON  запросы
	var task models.Task
	err := c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Вызов бизнес логики обновления задач
	code, err := h.Serv.UpdateTask(task, id)
	if err != nil {
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}
	c.Status(code)
}

// Обработчик логики удаления задач
func (h *Taskhandler) DeleteTaskHandler(c *gin.Context) {
	// Чтение параметров запроса
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id field is empty"})
		return
	}
	// Вызов бизнес логики удаления задачи
	code, err := h.Serv.DeleteTask(id)
	if err != nil {
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}
	c.Status(code)
}

// Обработчик логики завершения задачи
func (h *Taskhandler) FinishTaskHandler(c *gin.Context) {
	// Чтение параметров запроса
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id field is empty"})
		return
	}
	// Вызов сервисной функции
	code, err := h.Serv.ChangeStatus(id)
	if err != nil {
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}
	c.Status(code)
}

// Обработчик для получения информации о всех задачах
func (h *Taskhandler) TaskListsHandler(c *gin.Context) {
	// Чтение параметров запроса
	status := c.DefaultQuery("status", "active")
	if status != "active" && status != "done" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "status parameter is invalid"})
		return
	}
	// Вызов сервисной функции
	tasks, code, err := h.Serv.GetTasks(status)
	if err != nil {
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}
	c.JSON(code, tasks)
}

// Логика валидации задачи
func (h *Taskhandler) ValidateTask(task models.Task) error {
	if !task.Id.IsZero() {
		return errors.New("id string should be empty")
	}
	if task.ActiveDate == "" {
		return errors.New("active date field is not exist")
	}

	_, err := time.Parse("2006-01-02", task.ActiveDate)
	if err != nil {
		return errors.New("active date field is invalid: " + err.Error())
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
