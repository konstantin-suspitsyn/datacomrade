package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/konstantin-suspitsyn/datacomrade/query_engine/internal/config/constants"
	"github.com/konstantin-suspitsyn/datacomrade/query_engine/internal/db"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("godotenv file was not found")
	}

	dbConfig := constants.InitDbConfig()
	dbConnection, err := db.OpenDB(dbConfig.DB_USER, dbConfig.DB_PASSWORD, dbConfig.DB_HOST, dbConfig.DB_PORT, dbConfig.DB_DATABASE, dbConfig.DB_MAX_OPEN_CONNS, dbConfig.DB_MAX_IDLE_CONNS, dbConfig.DB_MAX_IDLE_TIME_MINS)
	if err != nil {
		log.Fatalf("failed to open database %q: %v\n", dbConfig.DB_DATABASE, err)
	}
	defer dbConnection.Close()

	log.Printf("✅ Connected to Postgres %q at %s:%d\n", dbConfig.DB_DATABASE, dbConfig.DB_HOST, dbConfig.DB_PORT)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down query_engine")
}
