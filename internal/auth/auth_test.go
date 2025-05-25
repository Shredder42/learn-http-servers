package auth

import (
	"testing"
)

func TestHashPasswordMatch(t *testing.T) {
	hash, err := HashPassword("baseball42")
	if err != nil {
		t.Error("couldn't hash password", err)
	}

	err = CheckPasswordHash(hash, "baseball42")
	if err != nil {
		t.Error("password does not match hash", err)
	}
}

func TestHashPasswordWrong(t *testing.T) {
	hash, err := HashPassword("doggyboy")
	if err != nil {
		t.Error("couldn't hash password", err)
	}

	err = CheckPasswordHash(hash, "doggyboyz")
	if err != nil {
		t.Error("password does not match hash", err)
	}
}
