package api

import (
	"github.com/yashbek/jotunheim/services/matchmaking"
	mainv1 "github.com/yashbek/y2j/api/main/v1"
)

type MainServer struct{
	mainv1.UnimplementedMainServiceServer
	MatchmakingPool *matchmaking.MatchmakingPool
}