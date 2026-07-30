package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/authlogicapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/tablesapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/userapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/userdomainrolesapiv1"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/config/constants"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/db"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/service/services"
	authlogicv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/auth_logic/v1"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	userv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user/v1"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Конфигурация сервера
// Настройка dev или prod
type APIServerConfiguration struct {
	Port int
	Env  string
}

func main() {

	if err := godotenv.Load(".env"); err != nil {
		panic("godotenv file was not found")
	}
	var serverConfig APIServerConfiguration
	hostPort := constants.InitHostPort()
	serverConfig.Port = hostPort.PORT

	flag.StringVar(&serverConfig.Env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// Потом надо переделать под разные env файлы
	envName := map[string]string{
		"development": ".env",
		"production":  ".env",
	}

	err := godotenv.Load(envName[serverConfig.Env])

	if err != nil {
		panic("godotenv file was not found")
	}

	//Инициализируем БД
	envConfig := constants.InitDbConfig()
	dbConnection, err := db.OpenDB(envConfig.DB_USER, envConfig.DB_PASSWORD, envConfig.DB_HOST, envConfig.DB_PORT, envConfig.DB_DATABASE, envConfig.DB_MAX_OPEN_CONNS, envConfig.DB_MAX_IDLE_CONNS, envConfig.DB_MAX_IDLE_TIME_MINS)
	if err != nil {
		log.Fatalf("failed to open database %q: %v\n", envConfig.DB_DATABASE, err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", hostPort.PORT))

	if err != nil {
		log.Printf("Failed to listen: %v\n", err)
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("Failed to close listner: %v\n", cerr)
		}
	}()

	// Сервисы
	seviceLayer := services.New(dbConnection)
	tablesapiv1 := tablesapiv1.New(seviceLayer)
	userapiv1 := userapiv1.New(seviceLayer)
	userdomainrolesapiv1 := userdomainrolesapiv1.New(seviceLayer)
	authlogicapiv1 := authlogicapiv1.New(seviceLayer)

	// Создаем GRPC сервер
	s := grpc.NewServer()

	// Регистрация сервисов
	userv1.RegisterUserServiceServer(s, userapiv1)
	tablesv1.RegisterTableServiceServer(s, tablesapiv1)
	userdomainrolesv1.RegisterUserDomainRolesServiceServer(s, userdomainrolesapiv1)
	authlogicv1.RegisterAuthLogicServiceServer(s, authlogicapiv1)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", serverConfig.Port)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
