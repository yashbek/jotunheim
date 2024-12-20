package api

import (
	"context"
	"fmt"

	mainv1 "github.com/yashbek/y2j/api/main/v1"
)

func (s MainServer) Ping(context.Context, *mainv1.PingRequest) (*mainv1.PingResponse, error) {
	resp := &mainv1.PingResponse{
		Html: "I EXIST TOO!!",
	}
	fmt.Print("HAY")

	return resp, nil
}
