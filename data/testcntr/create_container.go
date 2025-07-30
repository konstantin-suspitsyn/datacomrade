package testcntr

import (
	"context"
	"database/sql"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/konstantin-suspitsyn/datacomrade/configs"
	"github.com/konstantin-suspitsyn/datacomrade/db"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
	DB               *sql.DB
}

func New(ctx context.Context) (*PostgresContainer, error) {
	slog.Info("Starting Postgres container")
	postgresContainer, err := createPostgresContainer(ctx)
	if err != nil {
		return nil, err
	}
	slog.Info("Created Postgres container")
	connString, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	slog.Info(connString)

	DB, err := db.OpenDBWithConnString(connString)

	return &PostgresContainer{
		PostgresContainer: postgresContainer,
		ConnectionString:  connString,
		DB:                DB,
	}, nil

}
func createPostgresContainer(ctx context.Context) (*postgres.PostgresContainer, error) {

	err := godotenv.Load(filepath.Join("..", "..", ".dev-env"))

	if err != nil {
		panic("No env file")
	}

	envConfig := configs.InitTestDbConfig()

	return postgres.Run(ctx,
		envConfig.DB_CONTAINER_VERSION,
		postgres.WithInitScripts(filepath.Join("..", "..", "migrations", "schema.sql")),
		postgres.WithDatabase(envConfig.DB_DATABASE),
		postgres.WithUsername(envConfig.DB_USER),
		postgres.WithPassword(envConfig.DB_PASSWORD),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)

}
