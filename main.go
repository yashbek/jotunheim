package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"

	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/yashbek/jotunheim/api"
	"github.com/yashbek/jotunheim/api/httpapi"
	"github.com/yashbek/jotunheim/services/auth"
	"github.com/yashbek/jotunheim/websockets"

	"github.com/yashbek/jotunheim/db/firebasedb"
	mainv1 "github.com/yashbek/y2j/api/main/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load(".env")
	if err == nil {
		log.Println("reading from the .env file")
	}

	firebasedb.FirebaseClient, err = firebasedb.NewFirebaseApp(ctx, nil)
	if err != nil {
		log.Fatal("Couldn't initialize db", err)
	}

	mmQueue := api.Initalize()
	mainServer := api.MainServer{
		MatchmakingPool: &mmQueue,
		WSServer:        websockets.NewServer(),
	}

	go func() {
		if err := runRESTServer(ctx, mainServer); err != nil {
			log.Fatalf("Failed to run REST server: %v", err)
		}
	}()

	go func() {
		if err := runGRPCServer(ctx, mainServer); err != nil {
			log.Fatalf("Failed to run gRPC server: %v", err)
		}
	}()
	<-ctx.Done()
}

func loadTLSConfig(path string) *tls.Config {
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	caCert, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
		log.Fatal("Failed to append ca cert")
	}

	return &tls.Config{
		RootCAs: rootCAs,
		VerifyConnection: func(state tls.ConnectionState) error {
			opts := x509.VerifyOptions{
				DNSName: state.ServerName,
				Roots:   rootCAs,
			}
			_, err := state.PeerCertificates[0].Verify(opts)
			return err
		},
	}
}

func runGRPCServer(_ context.Context, server api.MainServer) error {
	const port = 8083
	const maxSizeInBytes = 1024 * 1024 * 8

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	opts = append(opts,
		grpc.MaxSendMsgSize(maxSizeInBytes),
		grpc.MaxRecvMsgSize(maxSizeInBytes),
		grpc.UnaryInterceptor(auth.AuthInterceptor),
	)
	grpcServer := grpc.NewServer(opts...)
	mainv1.RegisterMainServiceServer(grpcServer, server)

	reflection.Register(grpcServer)

	fmt.Printf("Starting gRPC server on :%d...\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

func runRESTServer(ctx context.Context, server api.MainServer) error {
	const port = 8082
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
		}),
		runtime.WithMetadata(auth.WithAuth),
	)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Connect-Protocol-Version", "Connect-Timeout-Ms", "X-Connect-Timeout-Ms", "authorization"},
		ExposedHeaders:   []string{"Content-Type", "Connect-Protocol-Version"},
		AllowCredentials: true,
	})

	mux.HandlePath("GET", "/", handleEntrance)
	mux.HandlePath("POST", "/signup", httpapi.SignupHandler)
	mux.HandlePath("POST", "/login", httpapi.LoginHandler)
	mux.HandlePath("GET", "/news", httpapi.HandleNews)
	mux.HandlePath("POST", "/matches", httpapi.MatchHistoryHandler)
	mux.HandlePath("GET", "/game", httpapi.GetGameHandler)
	mux.HandlePath("POST", "/evaluate", httpapi.EvaluateHandler)
	mux.HandlePath("GET", "/game/bot", httpapi.CreateBotGameHandler)

	creds := insecure.NewCredentials()
	conn, err := grpc.NewClient("localhost:8083",
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		return err
	}

	if err := mainv1.RegisterMainServiceHandler(ctx, mux, conn); err != nil {
		return err
	}

	mainHandler := http.NewServeMux()

	mainHandler.Handle("/", mux)

	mainHandler.HandleFunc("/ws", server.WSServer.HandleWS)

	fmt.Printf("Starting gRPC-Gateway server on :%d...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), c.Handler(mainHandler)); err != nil {
		return err
	}
	return nil
}

func handleEntrance(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.Write([]byte("pong"))
}
