package configs

type TestDBConfig struct {
	DB_DATABASE          string
	DB_USER              string
	DB_PASSWORD          string
	DB_CONTAINER_VERSION string
}

// Initializes config structure
func InitTestDbConfig() TestDBConfig {

	return TestDBConfig{
		DB_USER:              getStringEnv("POSTGRES_TEST_USER"),
		DB_PASSWORD:          getStringEnv("POSTGRES_TEST_PASSWORD"),
		DB_DATABASE:          getStringEnv("POSTGRES_TEST_DB"),
		DB_CONTAINER_VERSION: getStringEnv("POSTGRES_TEST_VERSION"),
	}
}
