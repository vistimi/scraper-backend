package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// DotEnvVariable use godot package to load/read the .env file and return the value of the key
func DotEnvVariable(key string) string {
	varEnv := os.Getenv(key)
	if varEnv == "" {
		log.Fatalf("Empty env variable %s", key)
	}
	return varEnv
}

func LoadEnvVariables(file string) {
	// local code for no containerization
	err := godotenv.Load(file)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
