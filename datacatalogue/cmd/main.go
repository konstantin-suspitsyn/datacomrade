package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	// Потом надо переделать под разные env файлы для разных окружений
	if err := godotenv.Load(".env"); err != nil {
		panic("godotenv file was not found")
	}

	var serverConfig APIServerConfiguration
	hostPort := constants.InitHostPort()
	serverConfig.Port = hostPort.PORT

	flag.StringVar(&serverConfig.Env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	//Инициализируем БД
	envConfig := constants.InitDbConfig()
	dbConnection, err := db.OpenDB(envConfig.DB_USER, envConfig.DB_PASSWORD, envConfig.DB_HOST, envConfig.DB_PORT, envConfig.DB_DATABASE, envConfig.DB_MAX_OPEN_CONNS, envConfig.DB_MAX_IDLE_CONNS, envConfig.DB_MAX_IDLE_TIME_MINS)
	if err != nil {
		log.Fatalf("failed to open database %q: %v\n", envConfig.DB_DATABASE, err)
	}
	defer dbConnection.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", hostPort.PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	// Сервисы
	serviceLayer := services.New(dbConnection)
	tablesAPI := tablesapiv1.New(serviceLayer)
	userAPI := userapiv1.New(serviceLayer)
	userDomainRolesAPI := userdomainrolesapiv1.New(serviceLayer)
	authLogicAPI := authlogicapiv1.New(serviceLayer)

	// Создаем GRPC сервер
	s := grpc.NewServer()

	// Регистрация сервисов
	userv1.RegisterUserServiceServer(s, userAPI)
	tablesv1.RegisterAliasServiceServer(s, tablesAPI)
	tablesv1.RegisterUserServiceServer(s, tablesAPI)
	tablesv1.RegisterHostServiceServer(s, tablesAPI)
	userdomainrolesv1.RegisterUserDomainRolesServiceServer(s, userDomainRolesAPI)
	authlogicv1.RegisterAuthLogicServiceServer(s, authLogicAPI)

	reflection.Register(s)

	serveErr := make(chan error, 1)
	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", serverConfig.Port)
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
		log.Println("🛑 Shutting down gRPC server...")
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

	log.Println("✅ Server stopped")
}
