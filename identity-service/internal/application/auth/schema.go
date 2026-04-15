package auth

type EmailVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenBody struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type RefreshRequest = TokenBody

type LogoutRequest = TokenBody

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
