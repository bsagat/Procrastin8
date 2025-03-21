package main

import (
	"log"
	"log/slog"

	"TodoApp/internal/app"
	repo "TodoApp/internal/dal"
)

func main() {
	slog.Info("Starting the program...")

	// Чтение конфигурации
	config := app.FetchConfig()
	slog.Info("Configuration parsing finished...")

	// Подключение к базе данных
	db, ctx, err := repo.Connect()
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Database connect finished...")

	defer db.Disconnect(ctx)

	// Создание роутера для сервера
	router := app.Setup(db)
	slog.Info("Server setup finished...")

	// Запуск сервера
	slog.Info("Server has been started on :" + *config.Port)
	if err := router.Run(":" + *config.Port); err != nil {
		log.Fatal("Server start error :" + err.Error())
	}
}
