package matchmaking

import (
	"slices"

	"github.com/yashbek/jotunheim/models"
	"github.com/yashbek/jotunheim/utils"
)

type MatchmakingPool []models.Profile

func Initalize() MatchmakingPool {
	return make([]models.Profile, 0)
}

func (pool *MatchmakingPool) NewPlayer(p models.Profile) {
	n := len((*pool))
	i, _ := slices.BinarySearchFunc((*pool), p, func(a, b models.Profile) int {
		return a.Elo - b.Elo
	})

	closestMatchIndex := i - 1

	if n == 0 {
		(*pool) = append((*pool), p)
		return
	}

	if utils.WithinBounds(i, 0, n) && !(utils.WithinBounds(i - 1, 0, n) && utils.Abs((*pool)[i - 1].Elo - p.Elo) < utils.Abs((*pool)[i].Elo - p.Elo)){
		closestMatchIndex = i
	}
	

	if utils.Abs((*pool)[closestMatchIndex].Elo - p.Elo) <= utils.DefaultEloInterval {
		closestMatch := (*pool)[closestMatchIndex]
		(*pool) = append((*pool)[:closestMatchIndex], (*pool)[closestMatchIndex+1:]...)
		launchMatch(p, closestMatch)
		return
	}

	switch {
	case i == 0:
		(*pool) = append([]models.Profile{p}, (*pool)...)
	case i == n:
		(*pool) = append((*pool), p)
	case utils.WithinBounds(i, 0, n):
		before := make(MatchmakingPool, len((*pool)[:i]))
		copy(before, (*pool)[:i])
		before = append(before, p)
		(*pool) = append(before, (*pool)[i:]...)
	}

}

func launchMatch(p1, p2 models.Profile) {

}
