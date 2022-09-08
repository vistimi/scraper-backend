package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// GetEnvVariable use godot package to load/read the .env file and return the value of the key
func GetEnvVariable(key string) string {
	varEnv := os.Getenv(key)
	if varEnv == "" {
		log.Fatalf("Error getting env variable %s", key)
	}
	return varEnv
}

func LoadEnvVariables(file string) {
	// local code for no containerization
	err := godotenv.Load(file)
	if err != nil {
		log.Fatalf("Error loading env file")
	}
}

func SetEnvVariable(key string, value string) {
	err := os.Setenv(key, value)
	if err != nil {
		log.Fatalf("Error setting env variable %s=%s", key, value)
	}
}

