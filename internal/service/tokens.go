package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secretKey = []byte("")

func SetSecretKey(key string) {
	secretKey = []byte(key)
}

func GenerateTokens(userID uuid.UUID) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
		"type":    "access",
	})

	accessTokenString, err := accessToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	// Refresh token также с ID пользователя
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(24 * 7 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"type":    "refresh",
	})

	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func ValidateRefreshToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Проверяем тип токена
		if claims["type"] != "refresh" {
			return uuid.Nil, errors.New("invalid token type")
		}

		// Получаем ID пользователя как строку
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return uuid.Nil, errors.New("invalid user_id format")
		}

		// Преобразуем строку в uuid.UUID
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return uuid.Nil, errors.New("invalid user_id")
		}

		return userID, nil
	}

	return uuid.Nil, errors.New("invalid token")
}
