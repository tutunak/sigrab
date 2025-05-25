package utils

import (
	"os"
	"path/filepath"
	"time"
)

var Timestamp string

func init() {
	// current timestamp
	Timestamp = time.Now().Format("20060102_150405")
}

func InitDir(path string) string {
	// Initialize the directory path with the current timestamp
	if path == "" {
		path = "/tmp/sigrab"
	}
	// create a directory with the current timestamp
	fullPath := filepath.Join(path, Timestamp)
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		panic("failed to create directory: " + err.Error())
	}
	return fullPath
}
