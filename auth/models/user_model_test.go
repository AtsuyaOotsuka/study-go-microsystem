package models

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestNewUser(t *testing.T) {
	name := "John Doe"
	email := "john.doe@example.com"
	password := "securepassword"
	refreshToken := "somerefreshtoken"

	user := NewUser(name, email, password, refreshToken)

	if user.Name != name {
		t.Errorf("expected name %q, got %q", name, user.Name)
	}

	if user.Email != email {
		t.Errorf("expected email %q, got %q", email, user.Email)
	}

	if user.Password != password {
		t.Errorf("expected password %q, got %q", password, user.Password)
	}

	if user.RefreshToken != refreshToken {
		t.Errorf("expected refresh token %q, got %q", refreshToken, user.RefreshToken)
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "securepassword"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	user := &User{
		Password: string(hashedPassword),
	}

	// 正しいパスワードで検証
	err = user.VerifyPassword(password)
	if err != nil {
		t.Errorf("expected password to be valid, got error: %v", err)
	}

	// 間違ったパスワードで検証
	err = user.VerifyPassword("wrongpassword")
	if err == nil {
		t.Error("expected error for invalid password, got nil")
	}
}
