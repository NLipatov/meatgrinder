package services

import (
	"fmt"
	"meatgrinder/internal/application/dtos"
)

type Command struct {
	Type        string
	CharacterID string
	Data        map[string]interface{}
}

func MapDTOToCommand(dto dtos.CommandDTO) (Command, error) {
	if dto.Type == "" {
		return Command{}, fmt.Errorf("command type is empty")
	}

	return Command{
		Type:        dto.Type,
		CharacterID: dto.CharacterID,
		Data:        dto.Data,
	}, nil
}
