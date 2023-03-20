package encryption

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		t.Errorf("Hashed password does not match original password: %v", err)
	}
}
