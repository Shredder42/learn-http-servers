package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPasswordMatch(t *testing.T) {
	wantErr := false
	hash, _ := HashPassword("baseball42")

	err := CheckPasswordHash(hash, "baseball42")
	if (err != nil) != wantErr {
		t.Error("password does not match hash", err)
	}
}

func TestHashPasswordWrong(t *testing.T) {
	wantErr := true
	hash, _ := HashPassword("doggyboy")

	err := CheckPasswordHash(hash, "doggyboyz")
	if (err != nil) != wantErr {
		t.Error("password does not match hash", err)
	}
}

func TestJWTCorrect(t *testing.T) {
	wantErr := false
	userID := uuid.New()
	tokenSecret := "ghost-is-the-best-jwt-key-ever-120723"

	tokenString, _ := MakeJWT(userID, tokenSecret, 1*time.Second)

	_, err := ValidateJWT(tokenString, tokenSecret)
	if (err != nil) != wantErr {
		t.Error("secret does not match", err)
	}
}

func TestJWTTimeOut(t *testing.T) {
	wantErr := true
	userID := uuid.New()
	tokenSecret := "ghost-is-the-best-jwt-key-ever-120723"

	tokenString, _ := MakeJWT(userID, tokenSecret, 1*time.Second)

	time.Sleep(2 * time.Second)

	_, err := ValidateJWT(tokenString, tokenSecret)
	if (err != nil) != wantErr {
		t.Error("JWT expired", err)
	}
}

func TestJWTSecretWrong(t *testing.T) {
	wantErr := true
	userID := uuid.New()
	tokenSecret := "ghost-is-the-best-jwt-key-ever-120723"
	otherSecret := "shredder-is-the-best-jwt-key-ever-120723"

	tokenString, _ := MakeJWT(userID, tokenSecret, 1*time.Second)

	_, err := ValidateJWT(tokenString, otherSecret)
	if (err != nil) != wantErr {
		t.Error("secret does not match", err)
	}
}
