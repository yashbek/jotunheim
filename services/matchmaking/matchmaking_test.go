package matchmaking

import (
	"slices"
	"testing"

	"github.com/yashbek/jotunheim/models"
)

func TestAdd(t *testing.T) {
	mmQueue := Initalize()
	newPlayers := []models.Profile{
		{
			ID:  "1",
			Elo: 0,
		},
		{
			ID:  "1.2",
			Elo: -100,
		},
		{
			ID:  "1.3",
			Elo: -101,
		},
		{
			ID:  "2",
			Elo: 400,
		},
		{
			ID:  "3",
			Elo: 100,
		},
		{
			ID:  "4",
			Elo: 300,
		},
		{
			ID:  "5",
			Elo: 420,
		},
	}

	for _, player := range newPlayers {
		mmQueue.newPlayer(player)
	}

	if !slices.IsSortedFunc(mmQueue, func(a, b models.Profile) int {
		return a.Elo - b.Elo
	}) {
		t.Errorf("expected")
	}

}
