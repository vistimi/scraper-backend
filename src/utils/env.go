package utils

import (
	"log"
	"os"
)

// DotEnvVariable use godot package to load/read the .env file and return the value of the key
func DotEnvVariable(key string) string {
	varEnv := os.Getenv(key)
	if varEnv == "" {
		log.Fatalf("Empty env variable %s", key)
	}
	return varEnv
}
