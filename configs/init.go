package configs

import (
	"fmt"
	"os"
	"strconv"
)

// Get any environment variable as string
// Panics if variable is not found
func getStringEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic(fmt.Sprintf("Env variable %s was not found", key))
}

// Get any environment variable as int64
// Panics if variable is not found or  van not be converter to int64
func getIntEnv(key string) int64 {
	value := getStringEnv(key)
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Env variable %s can not be converted to int64", key))
	}
	return intValue
}
