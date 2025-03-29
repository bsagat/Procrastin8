package app

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"TodoApp/internal/utils"
)

type SetupOptions struct {
	HelpFlag *bool
	Port     *string
}

// Инициализация конфигурации приложения
func FetchConfig() *SetupOptions {
	err := LoadFile("config/.env")
	if err != nil {
		log.Fatal(err)
	}

	opt := CheckFlags()
	return opt
}

// Логика чтения .env
func Parse(r io.Reader) (map[string]string, error) {
	envs := make(map[string]string)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Пропускаем комментарии и пустые строки
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Разбиваем строку на KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Убираем кавычки вокруг значений
		value = strings.Trim(value, `"'`)

		// Добавляем в мапу
		envs[key] = value
	}

	return envs, scanner.Err()
}

// Чтение .env файла
func LoadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err // Если файла нет, возвращаем ошибку
	}
	defer file.Close()

	envMap, err := Parse(file) // Читаем и парсим файл
	if err != nil {
		return err
	}
	// Устанавливаем переменные окружения
	for key, value := range envMap {
		os.Setenv(key, value)
	}
	return nil
}

// Получает информацию о флагах с командной строки
func CheckFlags() *SetupOptions {
	opts := &SetupOptions{}

	// Создаем флаги
	opts.HelpFlag = flag.Bool("help", false, "Prints information about program")
	opts.Port = flag.String("port", "8080", "Default port number")

	// Чтение командной строки
	flag.Parse()

	err := os.Setenv("port", *opts.Port)
	if err != nil {
		log.Fatal(err)
	}
	// Выводит информацию о запуске
	if *opts.HelpFlag {
		utils.PrintHelp()
	}
	return opts
}
