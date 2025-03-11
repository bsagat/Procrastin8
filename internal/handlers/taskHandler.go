// fmt.Println(c.Query("id"))
// fmt.Println(c.Params.ByName("id"))

// id, exist := c.Get("id")

// c.JSON(200, gin.H{"message": "pong", "id": id})

// fmt.Println(exist)
// fmt.Fprintln(c.Writer, c.Value("id"))
// fmt.Println(c.Keys)
// fmt.Println(c.Params.Get("gg"))

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

func (h *Taskhandler) NewTaskHandler(c *gin.Context) {
	var task models.Task
	err := c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = h.ValidateTask(task)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	code, err := h.Serv.CreateTask(&task)
	if err != nil {
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	c.JSON(code, gin.H{"inserted_id": task.Id})
}

func (h *Taskhandler) UpdateTaskHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id field is empty"})
		return
	}
	var task models.Task
	err := c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	code, err := h.Serv.UpdateTask(task, id)
	if err != nil {
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}
	c.Status(code)
}

func (h *Taskhandler) DeleteTaskHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id field is empty"})
		return
	}
	code, err := h.Serv.DeleteTask(id)
	if err != nil {
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}
	c.Status(code)
}

func (h *Taskhandler) FinishTaskHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id field is empty"})
		return
	}
	code, err := h.Serv.ChangeStatus(id)
	if err != nil {
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}
	c.Status(code)
}

func (h *Taskhandler) TaskListsHandler(c *gin.Context) {
	status := c.DefaultQuery("status", "active")
	tasks, code, err := h.Serv.GetTasks(status)
	if err != nil {
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}
	c.AsciiJSON(code, tasks)
}

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
		return errors.New("title field is not exis")
	}
	if len(task.Title) > 200 {
		return errors.New("title lenght is greater than 200")
	}
	return nil
}
