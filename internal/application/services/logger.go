package services

type Logger interface {
	LogEvent(event string)
}
