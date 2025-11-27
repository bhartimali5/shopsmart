package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

const JwtSecretKey = "secret-key"

func VerifyToken(tokenString string) (string, string, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(JwtSecretKey), nil
	})
	if err != nil {
		return "", "", errors.New("invalid token")
	}
	if !parsedToken.Valid {
		return "", "", errors.New("invalid token")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}
	userId := claims["user_id"].(string)
	userRole := claims["user_role"].(string)
	return userId, userRole, nil
}
