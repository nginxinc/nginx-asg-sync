package main

import "testing"

func TestSetPositiveIntOrDefault(t *testing.T) {
	defaultValue := 0
	value := setPositiveIntOrDefault(10, defaultValue)
	if value == 0 {
		t.Errorf(" setPositiveIntOrDefault() should return value %v but returned invalid value %v", value, defaultValue)
	}

	defaultValue = 1
	value = setPositiveIntOrDefault(0, defaultValue)
	if value != 1 {
		t.Errorf(" setPositiveIntOrDefault() should return default value %v but returned invalid value %v", defaultValue, value)
	}
}

func TestSetTimeOrDefault(t *testing.T) {
	defaultTime := "10s"
	time := setTimeOrDefault("20s", defaultTime)

	if time != "20s" {
		t.Errorf(" setTimeOrDefault() should return time %v but returned invalid time %v", time, defaultTime)
	}

	defaultTime = "10s"
	time = setTimeOrDefault("", defaultTime)

	if time != "10s" {
		t.Errorf(" setTimeOrDefault() should return default time %v but returned invalid time %v", defaultTime, time)
	}
}
