package util

import (
	"errors"
	"github.com/anish-yadav/lms-api/internal/constants"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

func CreateToken(data interface{}) (string, error) {
	secret := os.Getenv(constants.JwtSecret)
	if len(secret) == 0 {
		return "", errors.New(constants.JwtSecretNotFound)
	}
	claims := jwt.MapClaims{}
	claims["data"] = data
	claims["exp"] = time.Now().Add(time.Hour*24*31).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err

	}
	return token, err
}
