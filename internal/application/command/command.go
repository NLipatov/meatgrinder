package command

import (
	"fmt"
)

type Type int

const (
	UNSET Type = iota
	SPAWN
	MOVE
	ATTACK
	DISCONNECT
)

type Command struct {
	Type        Type
	CharacterID string
	Data        map[string]interface{}
}

func MapDTOToCommand(dto DTO) (Command, error) {
	if dto.Type == UNSET {
		return Command{}, fmt.Errorf("invalid command type is unset")
	}

	return Command{
		Type:        dto.Type,
		CharacterID: dto.CharacterID,
		Data:        dto.Data,
	}, nil
}
