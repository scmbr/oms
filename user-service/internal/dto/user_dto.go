package dto

import (
	"time"

	"github.com/scmbr/oms/user-service/internal/models"
)

type UserDTO struct {
	UserID           string    `json:"user_id"`
	Email            string    `json:"email"`
	Role             string    `json:"role"`
	RefreshToken     string    `json:"refresh_token"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

func ToUserDTO(user *models.User, refreshToken string, refreshExpiresAt time.Time) *UserDTO {
	return &UserDTO{
		UserID:           user.ID.String(),
		Email:            user.Email,
		Role:             user.Role,
		RefreshToken:     refreshToken,
		RefreshExpiresAt: refreshExpiresAt,
	}
}
