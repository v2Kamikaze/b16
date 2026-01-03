package domain

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, hashedPassword string) error
}

type TokenIssuer[ClaimsType any] interface {
	CreateToken(data ClaimsType) (string, error)
	DecodeToken(token string) (ClaimsType, error)
}
