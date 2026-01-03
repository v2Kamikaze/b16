package main

import (
	"fmt"
	"net/http"

	"github.com/v2code/b16/internal/auth/manager"
	"github.com/v2code/b16/internal/auth/policy"
	"github.com/v2code/b16/internal/config"
	"github.com/v2code/b16/internal/domain"
	"github.com/v2code/b16/internal/security"
)

func BasicAuthHandler(w http.ResponseWriter, r *http.Request, credentials domain.UserCredentials[*manager.BasicAuthCredentials]) {
	fmt.Fprintf(w, "Hello %v", credentials.Principal().Username)
}

func TokenAuthHandler(w http.ResponseWriter, r *http.Request, credentials domain.UserCredentials[*manager.TokenCredentials]) {
	fmt.Fprintf(w, "Hello %v", credentials.Principal().Email)
}

func main() {

	env := config.LoadEnvironment()

	jwtIssuer := security.NewJwtIssuer(env.TokenAuthEnv.Secret)

	basicAuthManager := manager.NewBasicAuthManager(manager.BasicAuthParams{
		Username: env.BasicAuthEnv.Username,
		Password: env.BasicAuthEnv.Password,
	})

	tokenAuthManager := manager.NewTokenAuthManager(jwtIssuer)

	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /basic-auth",
		domain.Auth(
			basicAuthManager,
			BasicAuthHandler,
		),
	)

	mux.HandleFunc(
		"GET /token-auth",
		domain.Auth(
			tokenAuthManager,
			domain.WithPolicy(
				TokenAuthHandler,
				policy.RequireRolePolicy("ADMIN", "USER"),
			),
		),
	)

	http.ListenAndServe(":8000", mux)
}
