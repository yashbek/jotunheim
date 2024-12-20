package api

import (
	"log"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/yashbek/jotunheim/db/firebasedb"
	"github.com/yashbek/jotunheim/models"
	services "github.com/yashbek/jotunheim/services/engine"
	"github.com/yashbek/jotunheim/utils"
	"github.com/yashbek/jotunheim/websockets"
	mainv1 "github.com/yashbek/y2j/api/main/v1"
)

type MainServer struct {
	mainv1.UnimplementedMainServiceServer
	MatchmakingPool *MatchmakingPool
	WSServer        *websockets.Server
}

type MatchmakingPool []models.User

func Initalize() MatchmakingPool {
	return make([]models.User, 0)
}

func (s MainServer) NewPlayer(p models.User) {
	pool := s.MatchmakingPool
	n := len((*pool))

	if slices.Contains(*pool, p) {
		return
	}

	i, _ := slices.BinarySearchFunc((*pool), p, func(a, b models.User) int {
		return a.Elo - b.Elo
	})

	closestMatchIndex := i - 1

	if n == 0 {
		(*pool) = append((*pool), p)
		return
	}

	if utils.WithinBounds(i, 0, n) && !(utils.WithinBounds(i-1, 0, n) && utils.Abs((*pool)[i-1].Elo-p.Elo) < utils.Abs((*pool)[i].Elo-p.Elo)) {
		closestMatchIndex = i
	}

	log.Print(*pool)

	if utils.Abs((*pool)[closestMatchIndex].Elo-p.Elo) <= utils.DefaultEloInterval {
		closestMatch := (*pool)[closestMatchIndex]
		(*pool) = append((*pool)[:closestMatchIndex], (*pool)[closestMatchIndex+1:]...)
		s.launchMatch(p, closestMatch)
		return
	}

	switch {
	case i == 0:
		(*pool) = append([]models.User{p}, (*pool)...)
	case i == n:
		(*pool) = append((*pool), p)
	case utils.WithinBounds(i, 0, n):
		before := make(MatchmakingPool, len((*pool)[:i]))
		copy(before, (*pool)[:i])
		before = append(before, p)
		(*pool) = append(before, (*pool)[i:]...)
	}

}

func (s MainServer) launchMatch(p1, p2 models.User) {
	gameID := uuid.New().String()

	newBoard := services.NewBoard(11)
	moves := services.MovesOrder{Moves: make([]services.Move, 0)}

	game := map[string]interface{}{
		"player1":    p1.Email,
		"player2":    p2.Email,
		"status":     "ready",
		"created_at": time.Now(),
		"board":      newBoard,
		"moves":      moves,
	}

	firebasedb.FirebaseClient.CreateGame("games", gameID, game)

	match1 := models.GameMatch{
		Opponent: p2,
		GameID:   gameID,
		Color:    "black",
	}
	match2 := models.GameMatch{
		Opponent: p1,
		GameID:   gameID,
		Color:    "white",
	}

	s.WSServer.SendMatchNotification(p1.Email, match1)
	s.WSServer.SendMatchNotification(p2.Email, match2)
}
