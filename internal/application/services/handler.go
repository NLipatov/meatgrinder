package services

import "meatgrinder/internal/application/command"

type Handler interface {
	Handle(c command.Command) error
}
