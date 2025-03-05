package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"TodoApp/internal/handlers"
	"TodoApp/internal/models"
	"TodoApp/internal/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	slog.Info("Starting the program...")

	err := utils.LoadFile(".env")
	if err != nil {
		log.Fatal(err)
	}

	utils.CheckFlags()

	client, err := Connect()
	if err != nil {
		slog.Error("Failed to Connect to the Database: " + err.Error())
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	new := models.Task{
		Id:         1,
		Title:      "title",
		ActiveDate: time.Now(),
		Status:     "active",
	}
	result, err := client.Database("testdb").Collection("users").InsertOne(context.TODO(), new)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	taskHandler := handlers.DefaultTaskHandler()

	router := gin.Default()

	router.POST("/api/todo-list/tasks", taskHandler.NewTaskHandler)
	router.PUT("/api/todo-list/tasks/:id", taskHandler.UpdateTaskHandler)
	router.DELETE("/api/todo-list/tasks/:id", taskHandler.DeleteTaskHandler)
	router.PUT("/api/todo-list/tasks/:id/done", taskHandler.FinishTaskHandler)
	router.GET("/api/todo-list/tasks", taskHandler.TaskListsHandler)

	slog.Info("Server has been started on :" + *models.Port)
	log.Fatal(router.Run(":" + *models.Port))
}

// Функция подключения к MongoDB
func Connect() (*mongo.Client, error) {
	// Получаем URI из переменной окружения
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, fmt.Errorf("переменная окружения MONGODB_URI не задана")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	fmt.Println(ctx)
	// Проверяем подключение
	if err := client.Ping(context.TODO(), nil); err != nil {
		client.Disconnect(context.TODO()) // Закрываем соединение при ошибке
		return nil, fmt.Errorf("MongoDB не отвечает: %w", err)
	}

	fmt.Println("✅ Успешное подключение к MongoDB")
	return client, nil
}
