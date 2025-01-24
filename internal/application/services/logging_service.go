package services

type LoggingService struct {
	logger Logger
}

func NewLoggingService(logger Logger) LoggingService {
	return LoggingService{
		logger: logger,
	}
}

func (s LoggingService) LogEvent(event string) {
	s.logger.LogEvent(event)
}
