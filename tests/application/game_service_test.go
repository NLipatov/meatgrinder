package application

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"math"
	"math/rand"
	"meatgrinder/internal/application/command"
	"meatgrinder/internal/application/services"
	"meatgrinder/internal/domain"
	"testing"
)

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(c command.Command) error {
	args := m.Called(c)
	return args.Error(0)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) LogEvent(event string) {
	m.Called(event)
}

func TestGameService_ProcessCommand(t *testing.T) {
	world := domain.NewWorld(1000, 1000)
	logger := new(MockLogger)
	logger.On("LogEvent", mock.AnythingOfType("string")).Return().Maybe()
	snapshotService := &services.WorldSnapshotService{}
	gameService := services.NewGameService(world, logger, snapshotService)

	t.Run("SPAWN command", func(t *testing.T) {
		charId := fmt.Sprintf("char-%v", rand.Intn(math.MaxInt32))
		cmd := command.Command{Type: command.SPAWN, CharacterID: charId}

		err := gameService.ProcessCommand(cmd)

		assert.NoError(t, err)
		assert.NotNil(t, world.Characters[charId])
	})

	t.Run("MOVE command", func(t *testing.T) {
		charId := fmt.Sprintf("char-%v", rand.Intn(math.MaxInt32))
		world.SpawnRandomCharacter(charId)
		character := world.Characters[charId]

		initialX, initialY := character.Position()
		cmd := command.Command{
			Type:        command.MOVE,
			CharacterID: charId,
			Data: map[string]interface{}{
				"dx": initialX + 1.0,
				"dy": initialY + 1.0,
			}}

		processCommandErr := gameService.ProcessCommand(cmd)
		assert.NoError(t, processCommandErr)

		world.Update()
		x, y := character.Position()

		assert.NotEqual(t, initialX, x)
		assert.NotEqual(t, initialY, y)
	})

	t.Run("ATTACK command", func(t *testing.T) {
		attackerId := fmt.Sprintf("char-%v", rand.Intn(math.MaxInt32))
		targetId := fmt.Sprintf("char-%v", rand.Intn(math.MaxInt32))

		mage := domain.NewMage(attackerId, 1, 1)
		warrior := domain.NewWarrior(targetId, 1, 1)
		world.Characters[attackerId] = mage
		world.Characters[targetId] = warrior

		_ = gameService.ProcessCommand(command.Command{
			Type:        command.MOVE,
			CharacterID: attackerId,
			Data: map[string]interface{}{
				"dx": 100.0,
				"dy": 100.0,
			}})

		targetInitH := world.Characters[targetId].Health()
		attackerInitH := world.Characters[attackerId].Health()

		_ = gameService.ProcessCommand(command.Command{
			Type:        command.ATTACK,
			CharacterID: attackerId,
			Data:        map[string]interface{}{"target_id": targetId},
		})

		gameService.UpdateWorld()

		assert.NotNil(t, world.Characters[attackerId])
		assert.NotNil(t, world.Characters[targetId])
		assert.Less(t, world.Characters[targetId].Health(), targetInitH)
		assert.Equal(t, world.Characters[attackerId].Health(), attackerInitH)
	})

	t.Run("DISCONNECT command", func(t *testing.T) {
		charId := fmt.Sprintf("char-%v", rand.Intn(math.MaxInt32))
		_ = gameService.ProcessCommand(command.Command{Type: command.SPAWN, CharacterID: charId})

		_, charWasOnMap := world.Characters[charId]

		_ = gameService.ProcessCommand(command.Command{Type: command.DISCONNECT, CharacterID: charId})
		_, charWasOnMapAfterDisconnect := world.Characters[charId]

		assert.True(t, charWasOnMap)
		assert.False(t, charWasOnMapAfterDisconnect)
	})

	t.Run("Unknown command", func(t *testing.T) {
		cmd := command.Command{Type: command.UNSET}

		err := gameService.ProcessCommand(cmd)

		assert.Error(t, err)
		assert.Equal(t, "unknown cmd 0", err.Error())
	})
}
