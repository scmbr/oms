package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey string
	ttl       time.Duration
}

func NewJWTManager(secretKey string, ttl time.Duration) *JWTManager {
	return &JWTManager{secretKey: secretKey, ttl: ttl}
}

func (j *JWTManager) Generate(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(j.ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTManager) Verify(tokenStr string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token.Claims.(*jwt.MapClaims), nil
}
