package auth

import (
	"testing"
)

func TestCheckHashEqual(t *testing.T) {
	password := "hunter2"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}
	if err = CheckPasswordHash(password, hash); err != nil {
		t.Fatalf("Password is not correct, should be")
	}
}

func TestCheckHashNotEqual(t *testing.T) {
	password := "hunter2"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}
	if err = CheckPasswordHash("password", hash); err == nil {
		t.Fatalf("Password is correct, shouldn't be")
	}
}
