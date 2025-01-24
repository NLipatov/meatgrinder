package services

type Handler interface {
	Handle(c Command) error
}
