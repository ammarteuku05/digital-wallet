package dto

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	FullName    string  `json:"full_name" validate:"required"`
	Email       *string `json:"email"`
	Password    string  `json:"password" validate:"required,min=6"`
	PhoneNumber string  `json:"phone_number" validate:"required,min=12"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresAt    string `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
