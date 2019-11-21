package main

import (
	"testing"
)

func TestGetFailTimeoutOrDefault(t *testing.T) {
	tests := []struct{
		input string
		expected string
	}{
		{
			input: "",
			expected: defaultFailTimeout,
		},
		{
			input: "10s",
			expected: "10s",
		},
	}

	for _, test := range tests {
		result := getFailTimeoutOrDefault(test.input)
		if result != test.expected {
			t.Errorf("getFailTimeoutOrDefault(%v) returned %v but expected %v", test.input, result, test.expected)
		}
	}
}

func TestGetSlowStartOrDefault(t *testing.T) {
	tests := []struct{
		input string
		expected string
	}{
		{
			input: "",
			expected: defaultSlowStart,
		},
		{
			input: "10s",
			expected: "10s",
		},
	}

	for _, test := range tests {
		result := getSlowStartOrDefault(test.input)
		if result != test.expected {
			t.Errorf("getSlowStartOrDefault(%v) returned %v but expected %v", test.input, result, test.expected)
		}
	}
}