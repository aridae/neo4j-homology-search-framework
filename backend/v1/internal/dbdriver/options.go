package dbdriver

import (
	"log"
	"os"
)

type Options struct {
	URI      string
	DB       string
	User     string
	Password string
}

func GetNeo4jOptions() *Options {
	dbUri, exists := os.LookupEnv("NEO4J_URI")
	if !exists {
		log.Fatal("database uri not found in .env")
	}

	dbName, exists := os.LookupEnv("NEO4J_DB")
	if !exists {
		log.Fatal("database username not found in .env")
	}

	dbUsr, exists := os.LookupEnv("NEO4J_USR")
	if !exists {
		log.Fatal("database username not found in .env")
	}

	dbPass, exists := os.LookupEnv("NEO4J_PASS")
	if !exists {
		log.Fatal("database password not found in .env")
	}

	return &Options{
		URI:      dbUri,
		DB:       dbName,
		User:     dbUsr,
		Password: dbPass,
	}
}
