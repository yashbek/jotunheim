package api

import (
	"context"

	"github.com/yashbek/jotunheim/models"
	mainv1 "github.com/yashbek/y2j/api/main/v1"
)

func (s MainServer) QueueUp(ctx context.Context, _ *mainv1.QueueUpRequest) (*mainv1.QueueUpResponse, error) {
	// @TODO: adding middleware to extract profile from context
	s.MatchmakingPool.NewPlayer(models.Profile{})

	resp := &mainv1.QueueUpResponse{
		Status: "InQueue",
	}

	return resp, nil
}
