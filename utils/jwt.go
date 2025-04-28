package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xyperam/go-spotify-clone/models"
)

var JWTKEY = []byte("inikeyy")

func GenerateJWT(user *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", user.ID),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTKEY)
	if err != nil {
		log.Println("Failed to sign token:", err)
		return "", err
	}
	fmt.Println("Token generated successfully")
	fmt.Println("Token:", tokenString)
	return tokenString, nil
}
