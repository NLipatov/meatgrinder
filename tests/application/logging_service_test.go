package application

import (
	"fmt"
	"math/rand"
	"meatgrinder/internal/application/services"
	"testing"
)

type mockLogger struct {
	logs []string
}

func (l *mockLogger) LogEvent(event string) {
	l.logs = append(l.logs, event)
}

func TestLoggingService_LogEvent(t *testing.T) {
	logger := &mockLogger{}
	ls := services.NewLoggingService(logger)

	logCount := rand.Intn(9) + 1
	for i := range logCount {
		log := fmt.Sprintf("event %d", i)
		ls.LogEvent(log)
	}

	if len(logger.logs) != logCount {
		t.Fatalf("Log count should be %d, got %d", logCount, len(logger.logs))
	}

	for i := range logger.logs {
		if logger.logs[i] != fmt.Sprintf("event %d", i) {
			t.Fatalf("Log event should be %s, got %s", fmt.Sprintf("event %d", i), logger.logs[i])
		}
	}
}
