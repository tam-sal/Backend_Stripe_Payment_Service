package env

import (
	"log/slog"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv(logger *slog.Logger) {
	err := godotenv.Load(".env.dev")
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace:", trace)
		os.Exit(1)
	}
}

func GetString(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}
	return value
}

func GetInt(key string, defaultvalue int) int {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultvalue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return intValue
}

func GetBool(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		panic(err)
	}
	return boolValue
}
