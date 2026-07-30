package constants

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// getStringEnv возвращает значение переменной окружения.
// Паникует, если переменная не найдена.
func getStringEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic(fmt.Sprintf("Env variable %s was not found", key))
}

// getIntEnv возвращает переменную окружения как int.
// Паникует, если переменная не найдена или не парсится в int.
func getIntEnv(key string) int {
	value := getStringEnv(key)
	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("Env variable %s can not be converted to int", key))
	}
	return intValue
}

// getBoolEnv возвращает переменную окружения как bool.
// Паникует, если переменная не найдена или не парсится в bool.
func getBoolEnv(key string) bool {
	value := getStringEnv(key)
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		panic(fmt.Sprintf("Env variable %s can not be converted to bool", key))
	}
	return boolValue
}

// getStringSliceEnv разбирает переменную окружения со списком значений через запятую.
// Паникует, если переменная не найдена.
func getStringSliceEnv(key string) []string {
	value := getStringEnv(key)

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}

	return result
}
