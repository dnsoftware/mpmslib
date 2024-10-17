package utils

import (
	"fmt"
	"os"
)

func SaveTextToFile(path, text string) error {
	// Создаем или открываем файл
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("ошибка при создании файла: %w", err)
	}
	defer file.Close() // Закрываем файл после записи

	// Записываем текст в файл
	_, err = file.WriteString(text)
	if err != nil {
		return fmt.Errorf("ошибка при записи текста в файл: %w", err)
	}

	return nil
}
