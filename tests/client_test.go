package tests

import (
	"context"
	"log"
	"testing"

	mainv1 "github.com/yashbek/y2j/api/main/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestPing (t *testing.T) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("localhost:8081", opts...)
	if err != nil {
		log.Fatal("grpc server connection failed")
	}
	defer conn.Close()

	client := mainv1.NewMainServiceClient(conn)

	req := &mainv1.PingRequest{}

	resp, err := client.Ping(context.Background(), req)
	if err != nil {
		log.Fatal("couldnt call Ping", err.Error())
	}
	t.Log(resp.Html)
}


