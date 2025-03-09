package app

import (
	repo "TodoApp/internal/dal"
	"TodoApp/internal/handlers"
	"TodoApp/internal/service"
	"TodoApp/internal/utils"
	"log"
	"log/slog"
)

// Вызывает функции которые нужны для запуска сервера (подключение к бд, чтение .env, командной строки)
func Setup() *handlers.Taskhandler {
	slog.Info("Starting the program...")

	err := utils.LoadFile(".env")
	if err != nil {
		log.Fatal(err)
	}

	utils.CheckFlags()

	db, err := repo.Connect()
	if err != nil {
		slog.Error("Failed to Connect to the Database: " + err.Error())
		log.Fatal(err)
	}

	taskRepo := repo.DefaultTaskRepository(db)
	taskService := service.DefaultTaskService(*taskRepo)
	taskHandler := handlers.DefaultTaskHandler(*taskService)
	return taskHandler
}
