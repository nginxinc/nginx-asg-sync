package main

import "testing"

func TestIsValidTime(t *testing.T) {
	testsWithValidInput := []string{"1", "1m10s", "11 11", "5m 30s", "1s", "100m", "5w", "15m", "11M", "3h", "100y", "600"}
	invalidInput := []string{"ss", "rM", "m0m", "s1s", "-5s", "1L"}

	for _, test := range testsWithValidInput {
		valid := isValidTime(test)
		if !valid {
			t.Errorf("isValidTime(%q) returned false for valid input.", test)
		}
	}
	for _, test := range invalidInput {
		valid := isValidTime(test)
		if valid {
			t.Errorf("isValidTime(%q) returned true for invalid input.", test)
		}
	}
}
