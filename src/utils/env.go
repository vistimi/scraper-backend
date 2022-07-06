package utils

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

// DotEnvVariable use godot package to load/read the .env file and return the value of the key
func DotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	varEnv := os.Getenv(key)
	if varEnv == "" {
		log.Fatalf("Empty env variable %s", key)
	}
	return varEnv
}
