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

func BasicAuthHandler(w http.ResponseWriter, r *http.Request, cred domain.UserCredentials[*manager.BasicAuthCredentials]) {
	fmt.Fprintf(w, "Hello %v", cred.GetCredentials().Username)
}

func TokenAuthHandler(w http.ResponseWriter, r *http.Request, cred domain.UserCredentials[*manager.TokenCredentials]) {
	fmt.Fprintf(w, "Hello %v", cred.GetCredentials().Email)
}

func main() {

	env := config.LoadEnvironment()

	jwtIssuer := security.NewJwtIssuer(env.TokenAuthEnv.Secret)

	token, err := jwtIssuer.CreateToken(&security.Claims{
		Email: "email@email.com",
		Roles: []string{"ADMIN"},
	})

	if err != nil {
		panic(err)
	}

	fmt.Print(token)

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
				policy.RequireRolePolicy("ADMIN"),
			),
		),
	)

	http.ListenAndServe(":8000", mux)
}
