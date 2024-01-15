package authentication

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const SecretKey = "superSecretKey"
const TokenExp = time.Hour * 24

type TokenError struct {
}

func (e TokenError) Error() string {
	return fmt.Sprintf("token is not valid")
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func BuildJWTString(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userId,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserIDFromJWTString(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		return "", TokenError{}
	}

	if !token.Valid {
		return "", TokenError{}
	}

	return claims.UserID, nil
}
