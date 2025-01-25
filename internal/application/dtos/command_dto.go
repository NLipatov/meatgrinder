package dtos

import "meatgrinder/internal/application/commands"

type CommandDTO struct {
	Type        commands.CommandType   `json:"type"`
	CharacterID string                 `json:"character_id"`
	Data        map[string]interface{} `json:"data"`
}
