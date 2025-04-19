package auth

type EmailVerificationRequestDTO struct {
	Email string `json:"email" binding:"required,email"`
}
