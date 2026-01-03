package security

import (
	"errors"

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
	jwt.RegisteredClaims
}

type JwtIssuer struct {
	secretKey     []byte
	signingMethod jwt.SigningMethod
}

func NewJwtIssuer(secretKey []byte) domain.TokenIssuer[*Claims] {
	return &JwtIssuer{
		secretKey:     secretKey,
		signingMethod: jwt.SigningMethodHS256,
	}
}

func (j *JwtIssuer) Create(claims *Claims) (string, error) {
	token := jwt.NewWithClaims(j.signingMethod, claims)

	signedToken, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *JwtIssuer) Decode(rawToken string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (any, error) {
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

	return claims, nil
}
