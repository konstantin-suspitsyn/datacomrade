package constants

type DbConfig struct {
	// Configurations for backend database
	DB_HOST               string
	DB_PORT               int
	DB_USER               string
	DB_PASSWORD           string
	DB_DATABASE           string
	DB_MAX_OPEN_CONNS     int
	DB_MAX_IDLE_CONNS     int
	DB_MAX_IDLE_TIME_MINS int
}

// Initializes config structure
func InitDbConfig() DbConfig {

	return DbConfig{
		DB_HOST:               getStringEnv("QUERY_ENGINE_POSTGRES_HOST"),
		DB_PORT:               int(getIntEnv("QUERY_ENGINE_POSTGRES_PORT")),
		DB_USER:               getStringEnv("QUERY_ENGINE_POSTGRES_USER"),
		DB_PASSWORD:           getStringEnv("QUERY_ENGINE_POSTGRES_PASSWORD"),
		DB_DATABASE:           getStringEnv("QUERY_ENGINE_POSTGRES_DB"),
		DB_MAX_OPEN_CONNS:     int(getIntEnv("DB_MAX_OPEN_CONNS")),
		DB_MAX_IDLE_CONNS:     int(getIntEnv("DB_MAX_IDLE_CONNS")),
		DB_MAX_IDLE_TIME_MINS: int(getIntEnv("DB_MAX_IDLE_TIME_MINS")),
	}
}
