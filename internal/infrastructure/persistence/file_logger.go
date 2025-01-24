package persistence

import (
	"fmt"
	"meatgrinder/internal/application/services"
	"os"
	"time"
)

type FileLogger struct {
	filePath string
}

func NewFileLogger(path string) services.Logger {
	return &FileLogger{filePath: path}
}

func (fl *FileLogger) LogEvent(event string) {
	f, err := os.OpenFile(fl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("FileLogger open error: %v\n", err)
		return
	}
	defer f.Close()

	line := fmt.Sprintf("%s | %s\n", time.Now().Format(time.RFC3339), event)
	if _, werr := f.WriteString(line); werr != nil {
		fmt.Printf("FileLogger write error: %v\n", werr)
	}
}
