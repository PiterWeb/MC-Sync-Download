package main

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

func getEnv(key string) (string, error) {

	envData, err := os.ReadFile(".env")

	if err != nil {
		return "", err
	}

	if key == "" {
		return "", errors.New("Key is empty")
	}

	env, err := godotenv.Unmarshal(string(envData))

	if err != nil {
		return "", err
	}

	return env[key], nil

}
