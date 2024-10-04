package db

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	maindb "github.com/yashbek/jotunheim/db/models/main"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBClientCtxKey struct{}

type DatabaseClient struct {
	*sqlx.DB
	isHealthy atomic.Bool
}

func Init(ctx context.Context, dbConfig DatabaseConfiguration) context.Context {
	client, ok := createDBConnection(dbConfig)

	if !ok {
		log.Fatal("failed to start db")
	}

	mainCtx := context.WithValue(ctx, DBClientCtxKey{}, client)

	return mainCtx
}

func GetDBClient(ctx context.Context) *DatabaseClient {
	client, ok := ctx.Value(DBClientCtxKey{}).(*DatabaseClient)

	if !ok {
		log.Fatal("context lacks client value, DB might have failed to initialize")
	}
	return client
}

func (dc *DatabaseClient) Query() *maindb.Queries {
	return maindb.New(dc.DB)
}

func createDBConnection(dbConfig DatabaseConfiguration) (*DatabaseClient, bool) {
	const maxRetries = 5

	host := dbConfig.Host
	port := dbConfig.Port
	username := dbConfig.User
	password := dbConfig.Password
	database := dbConfig.Database

	sqlxClient := initializeConnection(host, port, username, password, database, maxRetries)
	if sqlxClient != nil {
		dbClient := &DatabaseClient{
			sqlxClient,
			atomic.Bool{},
		}
		dbClient.isHealthy.Store(true)
		return dbClient, true
	}
	return &DatabaseClient{}, false
}

func initializeConnection(host string, port string, username string, password string, database string, retries int) *sqlx.DB {
	var err error
	var sqlxClient *sqlx.DB

	for sqlxClient == nil && retries > 0 {
		sqlxClient, err = sqlx.Connect("postgres", getConnectionString(host, port, username, password, database))
		retries--
	}
	for sqlxClient == nil {
		time.Sleep(time.Second)
		log.Println("Could not initialize mysql connection",
			err.Error(),
		)
		return nil
	}
	const defaultMaxOpenConns = 10
	const defaultMaxIdleConns = 5
	const defaultMaxConnAge = time.Minute * 3
	sqlxClient.SetMaxOpenConns(defaultMaxOpenConns)
	sqlxClient.SetMaxIdleConns(defaultMaxIdleConns)
	sqlxClient.SetConnMaxLifetime(defaultMaxConnAge)
	return sqlxClient
}

func getConnectionString(host string, port string, username string, password string, database string) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, database)
}

