package api

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/konstantin-suspitsyn/datacomrade/configs"
	"github.com/konstantin-suspitsyn/datacomrade/db"
	"github.com/konstantin-suspitsyn/datacomrade/internal/services"
	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/jsonlog"
)

type APIServerConfiguration struct {
	Port int
	Env  string
}

func StartServer() error {

	////////////////////////////////////////////////////////////
	////                      All Configs                   ////
	////////////////////////////////////////////////////////////

	var serverConfig APIServerConfiguration

	flag.IntVar(&serverConfig.Port, "port", 4000, "API Server Port")
	flag.StringVar(&serverConfig.Env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	envName := map[string]string{
		"development": ".dev-env",
		"production":  ".env",
	}

	err := godotenv.Load(envName[serverConfig.Env])

	if err != nil {
		panic("godotenv file was not found")
	}

	envConfig := configs.InitDbConfig()

	////////////////////////////////////////////////////////////
	////                      Shared resources              ////
	////////////////////////////////////////////////////////////

	// Pointer to a logger, that will be user everywhere in our program
	logger := createDefaultLogger(slog.LevelDebug)

	dbConnection, err := db.OpenDB(envConfig.DB_USER, envConfig.DB_PASSWORD, envConfig.DB_HOST, envConfig.DB_PORT, envConfig.DB_DATABASE, envConfig.DB_MAX_OPEN_CONNS, envConfig.DB_MAX_IDLE_CONNS, envConfig.DB_MAX_IDLE_TIME_MINS)

	if err != nil {
		panic(err)
	}

	serviceLayer := services.New(dbConnection)

	////////////////////////////////////////////////////////////
	////                      Router                        ////
	////////////////////////////////////////////////////////////

	router := routes(serviceLayer)
	////////////////////////////////////////////////////////////
	////                      Running server                ////
	////////////////////////////////////////////////////////////

	err = httpStart(router, logger, serverConfig)

	return err
}

// Start http server with graceful shutdown
func httpStart(mux *chi.Mux, loggerForServer *log.Logger, apiServerConfig APIServerConfiguration) error {

	var wg sync.WaitGroup

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", apiServerConfig.Port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * 60 * time.Second,
		WriteTimeout: 10 * 60 * time.Second,
		// Use custom logger
		ErrorLog: loggerForServer,
	}

	shutdownError := make(chan error)

	// Graceful shutdown handling
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		jsonlog.PrintInfo("Shutting down server", map[string]string{
			"signal": s.String(),
		}, nil)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)

		if err != nil {
			shutdownError <- err
		}

		jsonlog.PrintInfo("completing background tasks", map[string]string{
			"addr": srv.Addr,
		}, nil)

		wg.Wait()

		shutdownError <- nil
	}()

	jsonlog.PrintInfo("Started server", map[string]string{"addr": srv.Addr}, nil)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	jsonlog.PrintInfo("Stopped server", map[string]string{"addr": srv.Addr}, nil)

	return nil
}
