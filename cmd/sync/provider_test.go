package main

import "testing"

func TestValidateCloudProviderValid(t *testing.T) {
	t.Parallel()
	provider := "AWS"
	valid := validateCloudProvider(provider)
	if !valid {
		t.Errorf("validateCloudProvider(%v) returned invalid for a valid case", provider)
	}
}

func TestValidateCloudProviderInvalid(t *testing.T) {
	t.Parallel()
	provider := "invalid"
	valid := validateCloudProvider(provider)
	if valid {
		t.Errorf("validateCloudProvider(%v) returned valid for an invalid case", provider)
	}
}
