package services

type ILogger interface {
	LogEvent(event string)
}
