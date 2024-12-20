package api

import (
	"context"

	"github.com/yashbek/jotunheim/db/firebasedb"
	"github.com/yashbek/jotunheim/models"
	"github.com/yashbek/jotunheim/services/auth"
	mainv1 "github.com/yashbek/y2j/api/main/v1"
)

func (s MainServer) QueueUp(ctx context.Context, req *mainv1.QueueUpRequest) (*mainv1.QueueUpResponse, error) {
	claims := auth.FromCtx(ctx)
	email := claims.Email

	userInfo, err := firebasedb.FirebaseClient.Auth.GetUserByEmail(ctx, email)
	if err != nil {
		return &mainv1.QueueUpResponse{}, err
	}

	user := models.User{}
	firebasedb.FirebaseClient.ReadUser("users", userInfo.UID, &user)
	s.NewPlayer(user)

	resp := &mainv1.QueueUpResponse{
		Status: "InQueue",
	}

	return resp, nil
}
