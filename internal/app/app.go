package app

import (
	"net/http"
	"os"

	repo "TodoApp/internal/dal"
	"TodoApp/internal/handlers"
	"TodoApp/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Подключается к базе данных и создает роутер для приложения
func Setup(db *mongo.Client) *gin.Engine {
	taskRepo := repo.DefaultTaskRepository(db)
	taskService := service.DefaultTaskService(taskRepo)
	taskHandler := handlers.DefaultTaskHandler(taskService)

	router := gin.Default()
	router.LoadHTMLFiles("templates/index.html")

	router.GET("/api/todo-list", func(ctx *gin.Context) { // Основная страница
		domain := os.Getenv("API_DOMAIN")
		port := os.Getenv("port")
		url := domain + port + "/api/todo-list"
		ctx.HTML(http.StatusOK, "index.html", gin.H{"API_URL": url})
	})
	router.POST("/api/todo-list/tasks", taskHandler.NewTaskHandler)            // Создать одну задачу
	router.GET("/api/todo-list/tasks/:id", taskHandler.GetTaskHandler)         // Получить одну задачу из бд
	router.GET("/api/todo-list/tasks", taskHandler.TaskListsHandler)           // Получить список задач из бд
	router.PUT("/api/todo-list/tasks/:id", taskHandler.UpdateTaskHandler)      // Обновить существующую задачу
	router.DELETE("/api/todo-list/tasks/:id", taskHandler.DeleteTaskHandler)   // Удалить задачу
	router.PUT("/api/todo-list/tasks/:id/done", taskHandler.FinishTaskHandler) // Изменить статус задачи

	return router
}
