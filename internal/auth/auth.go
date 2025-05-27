package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	timeNow := jwt.NewNumericDate(time.Now().UTC())
	timeExpires := jwt.NewNumericDate(time.Now().UTC().Add(expiresIn))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  timeNow,
		ExpiresAt: timeExpires,
		Subject:   userID.String(),
	})

	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	if !token.Valid {
		return uuid.UUID{}, errors.New("invalid token")
	}

	// issuer, err := claims.GetIssuer() -- this also works
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.UUID{}, err
	}

	if issuer != "chirpy" {
		return uuid.UUID{}, errors.New("invalid issuer")
	}

	// userID, err := claims.GetSubject() -- this also works
	userID, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	return uuid.Parse(userID)

}
