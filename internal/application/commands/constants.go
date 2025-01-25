package commands

type CommandType int

const (
	UNSET CommandType = iota
	SPAWN
	MOVE
	ATTACK
	DISCONNECT
)
