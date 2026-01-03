package domain

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, hashedPassword string) error
}

type TokenIssuer[ClaimsT any] interface {
	Create(data ClaimsT) (string, error)
	Decode(token string) (ClaimsT, error)
}
