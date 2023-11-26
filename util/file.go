package util

import (
	"os"
)

func WriteFile(filePath string, content []byte) error {
	return os.WriteFile(filePath, content, 0o644)
}
