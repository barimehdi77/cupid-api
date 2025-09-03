package env

import (
	"os"
	"strconv"
)

func GetEnvString(key string, defaultValue string) string {
	env := os.Getenv(key)
	if env == "" {
		return defaultValue
	}
	return env
}

func GetEnvInt(key string, defaultValue int) int {
	env := GetEnvString(key, strconv.Itoa(defaultValue))
	port, _ := strconv.Atoi(env)
	return port
}
