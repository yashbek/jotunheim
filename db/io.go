package db

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/jmoiron/sqlx"
)



type DatabaseClient struct {
	*sqlx.DB
	isHealthy atomic.Bool
}

func LoadConfiguration(DatabaseConfiguration){

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

func createDBConnection(dbConfig DatabaseConfiguration) {
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
	}
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
		log.Print("Could not initialize mysql connection",
			err.Error(),
		)
		return nil
	}
	const defaultMaxOpenConns = 10
	const defaultMaxIdleConns = 5
	const readTime = time.Minute * 5
	const writeTime = time.Minute * 1
	defaultMaxConnAge := readTime
	sqlxClient.SetMaxOpenConns(defaultMaxOpenConns)
	sqlxClient.SetMaxIdleConns(defaultMaxIdleConns)
	sqlxClient.SetConnMaxLifetime(defaultMaxConnAge)
	return sqlxClient
}

func getConnectionString(host string, port string, username string, password string, database string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, database)
}