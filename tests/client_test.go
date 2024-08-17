package tests

// import (
// 	"context"
// 	"log"
// 	"testing"

// 	"github.com/yashbek/jotunheim.git/proto/generated"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// func TestPing (t *testing.T) {
// 	var opts []grpc.DialOption
// 	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	conn, err := grpc.NewClient("localhost:8081", opts...)
// 	if err != nil {
// 		log.Fatal("grpc server connection failed")
// 	}
// 	defer conn.Close()

// 	log.Println(conn)

// 	client := generated.NewMainClient(conn)

// 	req := &generated.PingRequest{}

// 	resp, err := client.Ping(context.Background(), req)
// 	if err != nil {
// 		log.Fatal("couldnt call Ping", err.Error())
// 	}
// 	log.Println(resp.Html)
// }


