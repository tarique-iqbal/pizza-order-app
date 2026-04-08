package auth

type OTPGenerator interface {
	Generate(secure bool) (string, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, plainPassword string) bool
}
