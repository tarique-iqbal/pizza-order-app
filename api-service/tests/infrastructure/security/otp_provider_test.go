package security_test

import (
	"api-service/internal/infrastructure/security"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate6DigitCode_Fast(t *testing.T) {
	otp := security.NewSixDigitOTPGenerator()
	for range 10 {
		code, err := otp.Generate(false)

		assert.NoError(t, err, "expected no error from fast code generator")
		assert.Len(t, code, 6, "code should be 6 characters long")
		assertDigitsOnly(t, code)
	}
}

func TestGenerate6DigitCode_Secure(t *testing.T) {
	otp := security.NewSixDigitOTPGenerator()
	for i := 0; i < 10; i++ {
		code, err := otp.Generate(true)

		assert.NoError(t, err, "expected no error from secure code generator")
		assert.Len(t, code, 6, "code should be 6 characters long")
		assertDigitsOnly(t, code)
	}
}

func assertDigitsOnly(t *testing.T, code string) {
	for _, ch := range code {
		assert.True(t, ch >= '0' && ch <= '9', "code should contain digits only")
	}
}
