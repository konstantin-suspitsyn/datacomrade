package testcntr

import (
	"context"
	"database/sql"
	"log/slog"
	"path/filepath"
	"runtime"
	"time"

	"github.com/konstantin-suspitsyn/datacomrade/configs"
	"github.com/konstantin-suspitsyn/datacomrade/db"
	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/comradetest"
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

	comradetest.InitEnv()

	envConfig := configs.InitTestDbConfig()
	_, dir, _, _ := runtime.Caller(0)

	env_dir := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(dir))))

	return postgres.Run(ctx,
		envConfig.DB_CONTAINER_VERSION,
		postgres.WithInitScripts(filepath.Join(env_dir, "migrations", "schema.sql")),
		postgres.WithDatabase(envConfig.DB_DATABASE),
		postgres.WithUsername(envConfig.DB_USER),
		postgres.WithPassword(envConfig.DB_PASSWORD),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)

}
