package utils

import (
	"os"
	"strconv"
)

func GetEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value != "" {
		val, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}
		return val
	}
	return defaultValue
}

func GetEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value != "" {
		val, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return val
	}
	return defaultValue
}

func GetEnvFloat(key string, defaultValue float64) float64 {
	value := os.Getenv(key)
	if value != "" {
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return defaultValue
		}
		return val
	}
	return defaultValue
}

func FirstEnv(keys ...string) (string, bool) {
	for _, k := range keys {
		if v, ok := os.LookupEnv(k); ok {
			return v, true
		}
	}
	return "", false
}

func Coalesce(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
