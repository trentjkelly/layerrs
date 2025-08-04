package config

import (
	"os"
	"fmt"

	"github.com/joho/godotenv"
)

const (
	ENVIRONMENT = "ENV"
)

// Gets the given environment for the application (this handles both local and docker deployments)
func GetEnvironment() (string, bool, error) {
	isDocker := false
	
	home, err := os.UserHomeDir()
	if err != nil {
		return "", false, fmt.Errorf("error getting user home directory: %w", err)
	}

	err = godotenv.Load(home + "/.env.layerrs")
	if err != nil {
		isDocker = true
	}

	env := os.Getenv(ENVIRONMENT)
	if env == "" {
		return "", false, fmt.Errorf("could not find the environment variable %s", ENVIRONMENT)
	}

	return env, isDocker, nil
}