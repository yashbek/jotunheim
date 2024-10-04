package db

import (
	"context"
	"log"
	"testing"
	"time"

	maindb "github.com/yashbek/jotunheim/db/models/main"

)

func TestPingDB(t *testing.T){
	ctx := Init(context.Background(), getDefaultDBConfig())

	client := GetDBClient(ctx)

	success := pingClient(client)

	if !success {
		t.Error("failed to ping DB")
	}

}

func pingClient(client *DatabaseClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 1)

	defer cancel()

	err := client.PingContext(ctx)

	if err != nil {
		log.Print(err.Error())
		return false
	}

	return true
}

func TestInsert(t *testing.T) {
	ctx := Init(context.Background(), getDefaultDBConfig())

	q := GetDBClient(ctx).Query()

	err := q.InsertProfile(ctx, maindb.InsertProfileParams{
		Username: "baloot",
		Email     : "@",
		PhoneNumber : "12",
		Elo : 12,
		DateJoined: time.Now(),
	})

	if err!=nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	ctx := Init(context.Background(), getDefaultDBConfig())

	q := GetDBClient(ctx).Query()

	profile, err := q.GetProfile(ctx, 5)
	if err!=nil {
		t.Error(err)
	}

	if profile.Elo != 12 {
		t.Error("profile mismatch")
	}
}