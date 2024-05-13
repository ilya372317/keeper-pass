package jwtmanager

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

type UserClaims struct {
	jwt.RegisteredClaims
	Email string
}

func New(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

func (j *JWTManager) Generate(user *domain.User) (string, error) {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
		},
		Email: user.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed sign jwt token: %w", err)
	}

	return tokenString, nil
}

func (j *JWTManager) Verify(accessToken string) (dto.JWTClaimsDTO, error) {
	token, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected token sign method")
		}

		return []byte(j.secretKey), nil
	})
	if err != nil {
		return dto.JWTClaimsDTO{}, fmt.Errorf("failed parse jwt token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return dto.JWTClaimsDTO{}, fmt.Errorf("unexcpeted token claims type")
	}

	return dto.JWTClaimsDTO{Email: claims.Email}, nil
}
