package configs

// import "github.com/joho/godotenv"

type JWTSecrets struct {
	JWT_SECRET string
}

func InitJwtSecrets() JWTSecrets {
	// godotenv.Load()
	return JWTSecrets{
		JWT_SECRET: getStringEnv("JWT_SECRET"),
	}
}
