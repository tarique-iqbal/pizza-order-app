package security

import (
	crand "crypto/rand"
	"fmt"
	"identity-service/internal/domain/auth"
	mrand "math/rand"
)

type SixDigitOTPGenerator struct {
}

func NewSixDigitOTPGenerator() auth.OTPGenerator {
	return &SixDigitOTPGenerator{}
}

func (g SixDigitOTPGenerator) Generate(secure bool) (string, error) {
	if secure {
		return generateSecureCode()
	}
	return generateFastCode()
}

func generateFastCode() (string, error) {
	n := mrand.Intn(1000000)
	return fmt.Sprintf("%06d", n), nil
}

func generateSecureCode() (string, error) {
	const max = 1000000
	var n int

	for {
		b := make([]byte, 4)
		if _, err := crand.Read(b); err != nil {
			return "", fmt.Errorf("crypto/rand failed: %w", err)
		}

		n = int(b[0])<<24 | int(b[1])<<16 | int(b[2])<<8 | int(b[3])

		if n >= 0 && n < (1<<32-(1<<32%max)) {
			n = n % max
			break
		}
	}

	return fmt.Sprintf("%06d", n), nil
}
