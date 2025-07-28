package configs

import "time"

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

type MailConfig struct {
	// Configuration for mail
	MAIL_HOST     string
	MAIL_PORT     int
	MAIL_USER     string
	MAIL_EMAIL    string
	MAIL_PASSWORD string
}

type JWTConfig struct {
	SecretKey            string
	TokenDuration        time.Duration
	RefreshTokenDuration time.Duration
}

func InitJWTConfig() JWTConfig {
	return JWTConfig{
		SecretKey:            getStringEnv("JWT_SECRET_KEY"),
		TokenDuration:        time.Duration(getIntEnv("JWT_DURATION_MINUTES")) * time.Minute,
		RefreshTokenDuration: time.Duration(getIntEnv("JWT_REFRESH_DURATION_IN_HRS")) * time.Hour,
	}
}

func InitMailConfig() MailConfig {
	// err := godotenv.Load()
	// if err != nil {
	// 	panic("Error loading .env file")
	// }

	return MailConfig{
		MAIL_HOST:     getStringEnv("MAIL_HOST"),
		MAIL_PORT:     int(getIntEnv("MAIL_PORT")),
		MAIL_USER:     getStringEnv("MAIL_USER"),
		MAIL_EMAIL:    getStringEnv("MAIL_EMAIL"),
		MAIL_PASSWORD: getStringEnv("MAIL_PASSWORD"),
	}
}

// Initializes config structure
func InitDbConfig() DbConfig {
	// err := godotenv.Load()
	// if err != nil {
	// 	panic("Error loading .env file")
	// }

	return DbConfig{
		DB_HOST:               getStringEnv("POSTGRES_BACKEND_HOST"),
		DB_PORT:               int(getIntEnv("POSTGRES_BACKEND_PORT")),
		DB_USER:               getStringEnv("POSTGRES_BACKEND_USER"),
		DB_PASSWORD:           getStringEnv("POSTGRES_BACKEND_PASSWORD"),
		DB_DATABASE:           getStringEnv("POSTGRES_BACKEND_DB"),
		DB_MAX_OPEN_CONNS:     int(getIntEnv("POSTGRES_BACKEND_MAX_OPEN_CONNECTIONS")),
		DB_MAX_IDLE_CONNS:     int(getIntEnv("POSTGRES_BACKEND_MAX_IDLE_CONNECTIONS")),
		DB_MAX_IDLE_TIME_MINS: int(getIntEnv("POSTGRES_CONN_MAX_IDLE_TIME_MINS")),
	}
}

