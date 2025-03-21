package app

import (
	repo "TodoApp/internal/dal"
	"TodoApp/internal/handlers"
	"TodoApp/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Подключается к базе данных и создает роутер для приложения
func Setup(db *mongo.Client) *gin.Engine {

	taskRepo := repo.DefaultTaskRepository(db)
	taskService := service.DefaultTaskService(*taskRepo)
	taskHandler := handlers.DefaultTaskHandler(*taskService)

	router := gin.Default()

	router.POST("/api/todo-list/tasks", taskHandler.NewTaskHandler)
	router.PUT("/api/todo-list/tasks/:id", taskHandler.UpdateTaskHandler)
	router.DELETE("/api/todo-list/tasks/:id", taskHandler.DeleteTaskHandler)
	router.PUT("/api/todo-list/tasks/:id/done", taskHandler.FinishTaskHandler)
	router.GET("/api/todo-list/tasks", taskHandler.TaskListsHandler)
	return router
}
