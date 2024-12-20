package tests

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHash(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	if err != nil {
		t.Error("idc")
	}

	newHash := string(hash)

	check := bcrypt.CompareHashAndPassword([]byte(newHash), []byte("12345678"))

	if check != nil {
		t.Error("problem")
	}

}
