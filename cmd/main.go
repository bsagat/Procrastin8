package main

import (
	"log"
	"log/slog"

	"TodoApp/internal/app"
	"TodoApp/internal/models"

	"github.com/gin-gonic/gin"
)

func main() {

	taskHandler := app.Setup()
	router := gin.Default()

	router.POST("/api/todo-list/tasks", taskHandler.NewTaskHandler)
	router.PUT("/api/todo-list/tasks/:id", taskHandler.UpdateTaskHandler)
	router.DELETE("/api/todo-list/tasks/:id", taskHandler.DeleteTaskHandler)
	router.PUT("/api/todo-list/tasks/:id/done", taskHandler.FinishTaskHandler)
	router.GET("/api/todo-list/tasks", taskHandler.TaskListsHandler)

	slog.Info("Server has been started on :" + *models.Port)
	log.Fatal(router.Run(":" + *models.Port))
}
