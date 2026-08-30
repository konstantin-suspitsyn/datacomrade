package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/konstantin-suspitsyn/datacomrade/query_engine/internal/api/queryengineapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/query_engine/internal/config/constants"
	"github.com/konstantin-suspitsyn/datacomrade/query_engine/internal/db"
	queryenginev1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/query_engine/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	hostPort := constants.InitHostPort()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", hostPort.PORT))
	if err != nil {
		log.Fatalf("failed to listen on %d: %v\n", hostPort.PORT, err)
	}

	s := grpc.NewServer()
	queryenginev1.RegisterQueryEngineServiceServer(s, queryengineapiv1.New())
	reflection.Register(s)

	serveErr := make(chan error, 1)
	go func() {
		log.Printf("🚀 query_engine gRPC server listening on %d\n", hostPort.PORT)
		if err := s.Serve(lis); err != nil {
			serveErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serveErr:
		log.Printf("gRPC server failed: %v\n", err)
	case <-quit:
		log.Println("🛑 Shutting down query_engine")
	}

	stopped := make(chan struct{})
	go func() {
		s.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
	case <-time.After(10 * time.Second):
		log.Println("graceful shutdown timed out, forcing stop")
		s.Stop()
	}

	log.Println("✅ query_engine stopped")
}
