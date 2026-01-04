package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/v2code/b16/internal/auth/manager"
	"github.com/v2code/b16/internal/auth/middleware"
	"github.com/v2code/b16/internal/auth/policy"
	"github.com/v2code/b16/internal/config"
	"github.com/v2code/b16/internal/domain"
	"github.com/v2code/b16/internal/security"
)

func BasicAuthHandler(w http.ResponseWriter, r *http.Request, credentials domain.Principal[*manager.BasicAuthPrincipal]) {
	fmt.Fprintf(w, "Hello %v", credentials.Principal().Username)
}

func TokenAuthHandler(w http.ResponseWriter, r *http.Request, credentials domain.Principal[*manager.TokenPrincipal]) {
	fmt.Fprintf(w, "Hello %v", credentials.Principal().Email)
}

func main() {

	env := config.LoadEnvironment()

	jwtIssuer := security.NewJwtIssuer(security.JwtIssuerParams{
		SecretKey: env.TokenAuthEnv.Secret,
		ExpireAt:  time.Hour * 2,
		Issuer:    "b16",
	})

	basicAuthManager := manager.NewBasicAuthManager(
		env.BasicAuthEnv.Username,
		env.BasicAuthEnv.Password,
	)

	tokenAuthManager := manager.NewTokenAuthManager(jwtIssuer)

	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /basic-auth",
		middleware.WithAuth(
			basicAuthManager,
			BasicAuthHandler,
		),
	)

	mux.HandleFunc(
		"GET /token-auth",
		middleware.WithAuth(
			tokenAuthManager,
			middleware.WithPolicy(
				TokenAuthHandler,
				policy.NewAnyPolicy(
					policy.RequireRolePolicy("ADMIN", "USER"),
					policy.NewCompositePolicy(
						policy.RequireRolePolicy("ADMIN", "USER"),
					),
				),
			),
		),
	)

	http.ListenAndServe(":8000", mux)
}
