package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/yashbek/jotunheim/api"
	"github.com/yashbek/jotunheim/services/matchmaking"
	mainv1 "github.com/yashbek/y2j/api/main/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	mmQueue := matchmaking.Initalize()
	mainServer := api.MainServer{
		MatchmakingPool: &mmQueue,
	}

	go func() {
		if err := runRESTServer(); err != nil {
			log.Fatalf("Failed to run REST server: %v", err)
		}
	}()

	go func() {
		if err := runGRPCServer(mainServer); err != nil {
			log.Fatalf("Failed to run gRPC server: %v", err)
		}
	}()	
}

func runGRPCServer(server api.MainServer) error {
	const port = 8081
	const maxSizeInBytes = 1024 * 1024 * 8
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	opts = append(opts, grpc.MaxSendMsgSize(maxSizeInBytes), grpc.MaxRecvMsgSize(maxSizeInBytes))
	grpcServer := grpc.NewServer(opts...)	
	mainv1.RegisterMainServiceServer(grpcServer, server)

	reflection.Register(grpcServer)

	fmt.Printf("Starting gRPC server on :%d...\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

func runRESTServer() error {
	const port = 8082
	ctx := context.Background()
	mux := runtime.NewServeMux()
	
	mux.HandlePath("GET", "/", handleEntrance)

	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	if err := mainv1.RegisterMainServiceHandler(ctx, mux, conn); err != nil {
		return err
	}

	fmt.Printf("Starting gRPC-Gateway server on :%d...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		return err
	}
	return nil
}

func handleEntrance (w http.ResponseWriter, _ *http.Request, _ map[string]string){
	w.Write([]byte("<h1>cheers!</h1>"))
}