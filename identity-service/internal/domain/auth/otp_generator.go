package auth

type OTPGenerator interface {
	Generate(secure bool) (string, error)
}
