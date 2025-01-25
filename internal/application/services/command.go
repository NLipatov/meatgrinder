package services

import (
	"fmt"
	"meatgrinder/internal/application/commands"
	"meatgrinder/internal/application/dtos"
)

type Command struct {
	Type        commands.CommandType
	CharacterID string
	Data        map[string]interface{}
}

func MapDTOToCommand(dto dtos.CommandDTO) (Command, error) {
	if dto.Type == commands.UNSET {
		return Command{}, fmt.Errorf("invalid command type is unset")
	}

	return Command{
		Type:        dto.Type,
		CharacterID: dto.CharacterID,
		Data:        dto.Data,
	}, nil
}
