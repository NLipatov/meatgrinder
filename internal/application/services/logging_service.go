package services

type ILogger interface {
	LogEvent(event string)
}

type LoggingService struct {
	logger ILogger
}

func NewLoggingService(logger ILogger) LoggingService {
	return LoggingService{
		logger: logger,
	}
}

func (s LoggingService) LogEvent(event string) {
	s.logger.LogEvent(event)
}
