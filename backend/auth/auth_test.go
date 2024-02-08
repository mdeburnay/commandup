package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "myPassword"

	hashedPassword, err := HashPassword(password)

	if err != nil {
		t.Error("Error hashing password")
	}

	if hashedPassword == password {
		t.Error("Password not hashed")
	}

	// Test for hash uniqueness
	hashedPassword2, err := HashPassword(password)
	if err != nil {
		t.Error("Error hashing password")
	}

	if hashedPassword == hashedPassword2 {
		t.Error("Hashed passwords are not unique")
	}

	// Verify the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		t.Error("Hashed password does not match the original password")
	}

	// Test hashing an empty string
	emptyHas, err := HashPassword("")
	if err != nil {
		t.Error("Error hashing empty string")
	}

	if emptyHas == "" {
		t.Error("Empty string not hashed")
	}
}
