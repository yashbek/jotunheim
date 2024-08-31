package api

import (
	"context"

	"github.com/yashbek/jotunheim/services"
	mainv1 "github.com/yashbek/y2j/api/main/v1"
)

func (s MainServer) QueueUp (ctx context.Context, _ *mainv1.QueueUpRequest) (*mainv1.QueueUpResponse, error) {
	services.AddToMatchMakingQueue("", 12)

	resp := &mainv1.QueueUpResponse{
		Status: "InQueue",
	}

	return resp, nil
}
