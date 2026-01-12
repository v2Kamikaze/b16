package config

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/v2code/b16/internal/logger"
)

type BasicAuth struct {
	Users map[string]string
}

type TokenAuth struct {
	Secret []byte
}

type Environment struct {
	TokenAuthEnv *TokenAuth
	BasicAuthEnv *BasicAuth
}

func LoadEnvironment() *Environment {
	err := godotenv.Load()
	if err != nil {
		logger.Debug("error while loading env: ", err.Error())
	}

	env := &Environment{}

	env.BasicAuthEnv = &BasicAuth{
		Users: LoadBasicAuthUsers(),
	}

	env.TokenAuthEnv = &TokenAuth{
		Secret: []byte(os.Getenv("B16_TOKEN_SECRET")),
	}

	return env
}

func ParseInt(value string) int {
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return parsedValue
}

func LoadBasicAuthUsers() map[string]string {
	usersJSON := os.Getenv("B16_BASIC_AUTH_USERS")
	if usersJSON == "" {
		panic("B16_BASIC_AUTH_USERS environment variable is required")
	}

	var users map[string]string
	if err := json.Unmarshal([]byte(usersJSON), &users); err != nil {
		panic("error parsing B16_BASIC_AUTH_USERS: " + err.Error())
	}

	return users
}
