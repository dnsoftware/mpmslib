package utils

import "os"

func FileExists(path string) bool {
	_, err := os.Stat(path)
	// os.IsNotExist проверяет, если ошибка означает отсутствие файла
	return !os.IsNotExist(err)
}
