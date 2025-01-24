package persistence

import (
	"fmt"
	"os"
	"time"
)

type FileLogger struct {
	filePath string
}

func NewFileLogger(path string) *FileLogger {
	return &FileLogger{filePath: path}
}

func (fl *FileLogger) LogEvent(event string) {
	f, err := os.OpenFile(fl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer func(f *os.File) {
		closeErr := f.Close()
		if closeErr != nil {
			fmt.Printf("Error closing file: %s\n", closeErr)
		}
	}(f)

	line := fmt.Sprintf("%s | %s\n", time.Now().Format(time.RFC3339), event)
	_, writeStringErr := f.WriteString(line)
	if writeStringErr != nil {
		fmt.Println("Error writing to file:", writeStringErr)
	}
}
