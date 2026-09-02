package config

import (
	"os"
	"strconv"
)

func Int(key string, fallback int) int {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	n, err := strconv.Atoi(value)

	if err != nil {
		return fallback
	}

	return n
}
