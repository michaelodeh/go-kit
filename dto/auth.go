package dto

import "github.com/golang-jwt/jwt/v5"

type LoginRequest struct {
	Email      string `json:"email" validate:"required"`
	Password   string `json:"password" validate:"required"`
	RememberMe bool   `json:"remember_me"`
}

type RegisterRequest struct {
	Name            string `json:"name" validate:"required,min=3,max=100"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6,max=100"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6,max=100"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ResetPasswordRequest struct {
	ID string `json:"id"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyEmailRequest struct {
	ID string `json:"id"`
}

type ChangePasswordRequest struct {
	Password        string `json:"password" validate:"required,min=6,max=100"`
	NewPassword     string `json:"new_password" validate:"required,min=6,max=100"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6,max=100"`
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
