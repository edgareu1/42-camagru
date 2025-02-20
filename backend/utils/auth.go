package utils

import (
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY = []byte("super-secret-key")

func GenerateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID,
		"ExpiresAt": time.Now().Add(24 * time.Hour),
	})
	return token.SignedString(SECRET_KEY)
}

func IsRequestAuthenticated(r *http.Request) bool {
	userIDStr := r.Header.Get("X-User-ID")
	authToken := r.Header.Get("X-Auth-Token")
	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})
	if err != nil {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	userID, err := strconv.ParseFloat(userIDStr, 64)
	if err != nil {
		return false
	}

	return claims["userID"] == userID
}
