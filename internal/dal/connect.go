package repo

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Функция подключения к MongoDB
func Connect() (*mongo.Client, error) {
	// Получаем URI из переменной окружения
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, fmt.Errorf("переменная окружения MONGODB_URI не задана")
	}

	_, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Проверяем подключение
	if err := client.Ping(context.TODO(), nil); err != nil {
		client.Disconnect(context.TODO()) // Закрываем соединение при ошибке
		return nil, fmt.Errorf("MongoDB не отвечает: %w", err)
	}

	fmt.Println("✅ Успешное подключение к MongoDB")
	return client, nil
}
