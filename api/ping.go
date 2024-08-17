package api

import (
	"context"

	mainv1 "github.com/yashbek/y2j/api/main/v1"
)

type MainServer struct{
	mainv1.UnimplementedMainServiceServer
}

func (s MainServer) Ping (context.Context, *mainv1.PingRequest) (*mainv1.PingResponse, error) {
	resp := &mainv1.PingResponse{
		Html: "<html> <h1>I EXIST TOO!!</h1> </html>",
	}
	
	return resp, nil
}
