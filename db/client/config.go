package db

import (
	"context"
	"log"
	"os"
	"strings"
)

const (
	defaultMainDB = "vchess"
	defaultMainDBHost = "localhost"
	defaultMainDBPort = "5432"
	defaultMainDBUser = "main"
	defaultMainDBPassword = "whocares"
)

type DatabaseConfiguration struct{
	Host string
	Port string
	User string
	Password string
	Database string
}

type DBCtxKey struct {}

func LoadDBConfig(ctx context.Context) (DatabaseConfiguration, context.Context) {
	dbConfig := getConfigFromEnv()

	return dbConfig, context.WithValue(ctx, DBCtxKey{}, dbConfig)
}

func getDefaultDBConfig() DatabaseConfiguration {
	return DatabaseConfiguration{
		Database: defaultMainDB,
		Host: defaultMainDBHost,
		Port: defaultMainDBPort,
		User: defaultMainDBUser,
		Password: defaultMainDBPassword,
	}
}

func getConfigFromEnv() DatabaseConfiguration {
	encodedConfigKey := "MAIN_DB" 
	configKeySeperator := ":"
	value := os.Getenv(encodedConfigKey)
	
	if value == "" {
		return getDefaultDBConfig()
	}

	values := strings.Split(value, configKeySeperator)

	if len(values) < 5 {
		log.Println("malformed db config encoded string, falling back to default")
		return getDefaultDBConfig()
	}

	dbConfig := DatabaseConfiguration {
		Database: values[0],
		Host: values[1],
		Port: values[2],
		User: values[3],
		Password: values[4],
	}

	return dbConfig
}