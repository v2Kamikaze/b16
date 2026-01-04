package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/v2code/b16/internal/domain"
)

var (
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidToken         = errors.New("invalid token")
)

type Claims struct {
	Email string
	Roles []string
}

type internalClaims struct {
	Email string
	Roles []string
	jwt.RegisteredClaims
}

type JwtIssuer struct {
	secretKey     []byte
	expireAt      time.Duration
	issuer        string
	signingMethod jwt.SigningMethod
}

type JwtIssuerParams struct {
	SecretKey []byte
	ExpireAt  time.Duration
	Issuer    string
}

func NewJwtIssuer(params JwtIssuerParams) domain.TokenIssuer[*Claims] {
	return &JwtIssuer{
		secretKey:     params.SecretKey,
		expireAt:      params.ExpireAt,
		issuer:        params.Issuer,
		signingMethod: jwt.SigningMethodHS256,
	}
}

func (j *JwtIssuer) Create(claims *Claims) (string, error) {

	token := jwt.NewWithClaims(j.signingMethod, internalClaims{
		Roles: claims.Roles,
		Email: claims.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expireAt)),
		},
	})

	signedToken, err := token.SignedString(j.secretKey)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *JwtIssuer) Decode(rawToken string) (*Claims, error) {

	internalClaims := &internalClaims{}

	token, err := jwt.ParseWithClaims(rawToken, internalClaims, func(token *jwt.Token) (any, error) {
		if token.Method != j.signingMethod {
			return nil, ErrInvalidSigningMethod
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return &Claims{
		Email: internalClaims.Email,
		Roles: internalClaims.Roles,
	}, nil
}
