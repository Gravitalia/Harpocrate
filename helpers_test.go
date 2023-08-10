package main

import (
	"testing"

	"github.com/Gravitalia/Harpocrate/helpers"
)

// TestGenerateRandomString calls helpers.GenerateRandomString with a
// specified lenght, checking for a valid return value.
func TestGenerateRandomString(t *testing.T) {
	msg, err := helpers.GenerateRandomString(8)
	if len(msg) != 8 || err != nil {
		t.Fatalf(`GenerateRandomString("") = %q, %v, want "", error`, msg, err)
	}
}
