package infrastructure

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "super-secret-key"

func TestGenerateJWT(t *testing.T) {
	token, err := GenerateJWT(
		1,
		"admin@example.com",
		testSecret,
		time.Hour,
	)

	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}

	if token == "" {
		t.Fatal("expected token, got empty string")
	}
}

func TestParseJWT(t *testing.T) {
	token, err := GenerateJWT(
		1,
		"admin@example.com",
		testSecret,
		time.Hour,
	)

	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}

	claims, err := ParseJWT(token, testSecret)
	if err != nil {
		t.Fatalf("ParseJWT() error = %v", err)
	}

	if claims.UserID != 1 {
		t.Errorf("expected UserID=1, got %d", claims.UserID)
	}

	if claims.Email != "admin@example.com" {
		t.Errorf("expected Email=admin@example.com, got %s", claims.Email)
	}

	if claims.Subject != "admin@example.com" {
		t.Errorf("expected Subject=admin@example.com, got %s", claims.Subject)
	}
}

func TestParseJWT_InvalidSecret(t *testing.T) {
	token, err := GenerateJWT(
		1,
		"admin@example.com",
		testSecret,
		time.Hour,
	)

	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}

	_, err = ParseJWT(token, "wrong-secret")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseJWT_ExpiredToken(t *testing.T) {
	token, err := GenerateJWT(
		1,
		"admin@example.com",
		testSecret,
		-time.Hour,
	)

	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}

	_, err = ParseJWT(token, testSecret)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Fatalf("expected ErrTokenExpired, got %v", err)
	}
}

func TestParseJWT_InvalidToken(t *testing.T) {
	_, err := ParseJWT("this-is-not-a-jwt", testSecret)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
