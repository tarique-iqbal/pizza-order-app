package auth

import "errors"

var (
	ErrCodeInvalid   = errors.New("invalid verification code")
	ErrCodeExpired   = errors.New("verification code expired")
	ErrCodeUsed      = errors.New("verification code already used")
	ErrCodeNotIssued = errors.New("verification code not issued")
)
