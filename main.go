package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/v2code/b16/internal/auth"
	"github.com/v2code/b16/internal/auth/manager"
	"github.com/v2code/b16/internal/auth/middleware"
	"github.com/v2code/b16/internal/auth/policy"
	"github.com/v2code/b16/internal/config"
	"github.com/v2code/b16/internal/logger"
	"github.com/v2code/b16/internal/mailer"
	"github.com/v2code/b16/internal/security"
)

func BasicAuthHandler(w http.ResponseWriter, r *http.Request, principal auth.Principal[*manager.BasicAuthPrincipal]) {
	fmt.Fprintf(w, "Hello %v\n", principal.Principal().Username)
}

func TokenAuthHandler(w http.ResponseWriter, r *http.Request, principal auth.Principal[*manager.TokenPrincipal]) {
	fmt.Fprintf(w, "Hello %v\n", principal.Principal().Email)
}

func main() {

	env := config.LoadEnvironment()

	jwtIssuer := security.NewJwtIssuer(security.JwtIssuerParams{
		SecretKey: env.TokenAuthEnv.Secret,
		ExpireAt:  time.Hour * 2,
		Issuer:    "b16",
	})

	basicAuthManager := manager.NewBasicAuthManager(env.BasicAuthEnv.Users)

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

	logger.Debug("server is running", "url", "http://0.0.0.0:8000")

	c := mailer.NewDefaultClient()

	m := mailer.NewDefaultMailer(mailer.MailerParams{
		Host:     "localhost",
		Port:     1025,
		From:     "b16@email.com",
		Username: "",
		Password: "",
	}, c)

	body := mailer.RenderVerificationCodeTemplate("123456")

	err := m.Send("Verification Code", body, "recipient@example.com")
	if err != nil {
		logger.Error("failed to send email", "error", err)
	}

	http.ListenAndServe(":8000", mux)
}
