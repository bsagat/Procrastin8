package utils

import (
	"TodoApp/internal/models"
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func CheckFlags() {
	flag.Parse()
	if *models.HelpFlag {
		PrintHelp()
	}
}

func PrintHelp() {
	fmt.Println(models.HelpInformation)
	os.Exit(0)
}

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
