package auth

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestCheckJWTValid(t *testing.T) {
	id := uuid.New()
	token, err := MakeJWT(id, "secret", 5*time.Second)
	if err != nil {
		t.Fatalf("Error generating token: %v", err)
	}
	returnedID, err := ValidateJWT(token, "secret")
	if err != nil {
		t.Fatalf("Error validating token: %v", err)
	}
	if id != returnedID {
		t.Fatalf("Expected IDs to match")
	}
}

func TestCheckJWTInvalidExpired(t *testing.T) {
	id := uuid.New()
	token, err := MakeJWT(id, "secret", 5*time.Millisecond)
	if err != nil {
		t.Fatalf("Error generating token; %v", err)
	}
	time.Sleep(10 * time.Millisecond)
	_, err = ValidateJWT(token, "secret")
	if err == nil {
		t.Fatalf("Expected expired token")
	} else if !strings.Contains(err.Error(), "token is expired") {
		t.Fatalf("Expected error: token is expired, got %v", err)
	}
}

func TestCheckJWTInvalidWrongSecret(t *testing.T) {
	id := uuid.New()
	token, err := MakeJWT(id, "secret", 5*time.Second)
	if err != nil {
		t.Fatalf("Error generating token: %v", err)
	}
	_, err = ValidateJWT(token, "terces")
	if err == nil {
		t.Fatalf("Expected expired token")
	} else if !strings.Contains(err.Error(), "signature is invalid") {
		t.Fatalf("Expected error: signature is invalid, got %v", err)
	}
}

func TestCheckBearerValid(t *testing.T) {
	header := http.Header{}
	header.Add("Authorization", "Bearer tokentest")
	token, err := GetBearerToken(header)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token != "tokentest" {
		t.Fatalf("expected 'tokentest', got %s", token)
	}
}

func TestCheckBearerNoBearer(t *testing.T) {
	header := http.Header{}
	token, err := GetBearerToken(header)
	if err == nil {
		t.Fatalf("expected no return, got %v", token)
	}
}
