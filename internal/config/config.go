package config

import (
	"os"
	"strconv"
)

type BasicAuth struct {
	Username string
	Password string
}

type TokenAuth struct {
	Secret []byte
}

type Environment struct {
	TokenAuthEnv *TokenAuth
	BasicAuthEnv *BasicAuth
}

func LoadEnvironment() *Environment {
	env := &Environment{}

	env.BasicAuthEnv = &BasicAuth{
		Username: GetEnvOrDefault("B16_USERNAME", "admin"),
		Password: GetEnvOrDefault("B16_PASSWORD", "password"),
	}

	env.TokenAuthEnv = &TokenAuth{
		Secret: []byte(GetEnvOrDefault("B16_TOKEN_SECRET", "secret")),
	}

	return env
}

func GetEnvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func ParseInt(value string) int {
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		panic(err.Error())
	}
	return parsedValue
}
