package auth

type CodeVerifier interface {
	Verify(email string, code string) error
}
