package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHashingPassword(t *testing.T) {
	password := "MatiasPassword"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Errorf(`HashPassword Fail %q, %v`, hashed, err)
	}
}

// TestHelloEmpty calls greetings.Hello with an empty string,
// checking for an error.
func TestCheckingPassword(t *testing.T) {
	password := "MatiasPassword"
	hashed, _ := HashPassword(password)
	err := CheckPasswordHash(password, hashed)
	if err != nil {
		t.Errorf(`Checking password fail %v, want "", error`, err)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJwtFuncs(t *testing.T) {
	tokenSecret := "*SXDI5!uv9akdS6u&s8gAAJBDJAsdAOHOLKNZx&$c0"
	userID := uuid.MustParse("1025b229-8fcb-4b25-81b9-ada154d526e6")
	var duration time.Duration = 1000000000 * 5
	jwtUser, err := MakeJWT(userID, tokenSecret, duration)
	if err != nil {
		t.Errorf(`fail make jwt %v, want "", error`, err)
	}

	validateId, err := ValidateJWT(jwtUser, tokenSecret)
	if err != nil {
		t.Errorf(`fail validate jwt %v, want "", error`, err)
	}
	if validateId != userID {
		t.Errorf(`is not the same uuid pass validation  %v, want %v, error`, validateId, userID)
	}

}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}
