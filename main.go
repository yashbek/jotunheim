package main

import (
	"fmt"
	"log"
	"net"

	mainv1 "github.com/yashbek/y2j/api/main/v1"
	"github.com/yashbek/jotunheim/api"
	"google.golang.org/grpc"
)

func main() {
	port := 8081
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
	log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	mainv1.RegisterMainServiceServer(grpcServer, api.MainServer{})
	grpcServer.Serve(lis)
	
}
