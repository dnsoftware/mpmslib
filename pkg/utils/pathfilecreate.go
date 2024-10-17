package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateFileWithDirs(path string) error {
	// Извлекаем путь к директории
	dir := filepath.Dir(path)

	// Создаем все директории, если их нет
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("ошибка при создании директорий: %w", err)
	}

	// Создаем файл
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("ошибка при создании файла: %w", err)
	}
	defer file.Close() // Закрываем файл после использования

	return nil
}
