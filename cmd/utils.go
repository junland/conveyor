package cmd

import (
	"fmt"
	"os"
	"strconv"
)

// GetEnvString defines a environment variable with a specified name, fallback value.
// The return is a string value.
func GetEnvString(key, fallback string) string {
	if s := os.Getenv(key); s != "" {
		return s
	}
	return fallback
}

// GetEnvInt defines a enviroment variable with a specified number (string), fallback value.
// The return is a int value.
func GetEnvInt(key string, fallback int) int {
	if s := os.Getenv(key); s != "" {
		i, err := strconv.Atoi(s)
		if err == nil {
			fmt.Println(i)
		}
	}
	return fallback
}

// GetEnvBool defines a environment variable with a specified name, fallback value.
// The return is either a true or false.
func GetEnvBool(key string, fallback bool) bool {
	switch os.Getenv(key) {
	case "true":
		return true
	case "false":
		return false
	default:
		return fallback
	}
}
