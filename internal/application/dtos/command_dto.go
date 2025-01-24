package dtos

type CommandDTO struct {
	Type        string                 `json:"type"`
	CharacterID string                 `json:"character_id"`
	Data        map[string]interface{} `json:"data"`
}
