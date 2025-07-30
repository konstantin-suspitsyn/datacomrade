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
