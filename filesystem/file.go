package filesystem

import (
	"errors"
	"os"
)

func FileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}
