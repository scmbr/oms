package service

import "time"

type UserDTO struct {
	UserID           string    `json:"user_id"`
	Email            string    `json:"email"`
	Role             string    `json:"role"`
	RefreshToken     string    `json:"refresh_token"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

func toUserDTO(userID, email, role, refreshToken string, refreshExpiresAt time.Time) *UserDTO {
	return &UserDTO{
		UserID:           userID,
		Email:            email,
		Role:             role,
		RefreshToken:     refreshToken,
		RefreshExpiresAt: refreshExpiresAt,
	}
}
