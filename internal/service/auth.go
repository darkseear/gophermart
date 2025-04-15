package service

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Auth struct {
	secretKey string
}

func NewAuth(secretKey string) *Auth {
	return &Auth{secretKey: secretKey}
}

func (s *Auth) GenerateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *Auth) ValidateToken(token string) (int, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		userID := claims["userID"].(float64)
		return int(userID), nil
	}
	return 0, errors.New("invalid token")
}

func (s *Auth) RefreshToken(token string) (string, error) {
	userID, err := s.ValidateToken(token)
	if err != nil {
		return "", err
	}
	return s.GenerateToken(userID)
}

func (s *Auth) ExtractToken(token string) (int, error) {
	userID, err := s.ValidateToken(token)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
