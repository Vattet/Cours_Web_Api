package service

import (
	"os"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)

func Login(username string, password string) (string, error) {
	user := os.Getenv("ADMIN_USERNAME")
	passwd := os.Getenv("ADMIN_PASSWORD")

	if user != username {
		return "", errors.New("Wrong Identifiants")
	}

	if bcrypt.CompareHashAndPassword([]byte(passwd), []byte(password)) != nil {
		return "", errors.New("Wrong Identifiants")
	}

	claims := jwt.RegisteredClaims{Subject: username, IssuedAt: jwt.NewNumericDate(time.Now()), ExpiresAt: jwt.NewNumericDate(time.Now().Add(24*time.Hour))}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", errors.New("Can't singed claims JWT")
	}

	return signed, nil
}

