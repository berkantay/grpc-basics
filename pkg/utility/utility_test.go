package utility

import (
	"testing"
)

func TestCheckIsValidMail(t *testing.T) {
	validEmail := "test@example.com"
	if !CheckIsValidMail(validEmail) {
		t.Errorf("Expected '%s' to be a valid email address", validEmail)
	}

	invalidEmail := "example.com"
	if CheckIsValidMail(invalidEmail) {
		t.Errorf("Expected '%s' to be an invalid email address", invalidEmail)
	}
}
