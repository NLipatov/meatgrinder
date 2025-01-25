package command

type DTO struct {
	Type        Type                   `json:"type"`
	CharacterID string                 `json:"character_id"`
	Data        map[string]interface{} `json:"data"`
}
