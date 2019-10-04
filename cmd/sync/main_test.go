package main

import "testing"

func TestSetPositiveInt(t *testing.T) {
	defaultValue := 0
	value := setPositiveInt(10, defaultValue)
	if value == 0 {
		t.Errorf(" setPositiveInt() should return value %v but returned invalid value %v", value, defaultValue)
	}

	defaultValue = 1
	value = setPositiveInt(0, defaultValue)
	if value != 1 {
		t.Errorf(" setPositiveInt() should return default value %v but returned invalid value %v", defaultValue, value)
	}
}

func TestSetTime(t *testing.T) {
	defaultTime := "10s"
	time := setTime("20s", defaultTime)

	if time != "20s" {
		t.Errorf(" setTime() should return time %v but returned invalid time %v", time, defaultTime)
	}

	defaultTime = "10s"
	time = setTime("", defaultTime)

	if time != "10s" {
		t.Errorf(" setTime() should return default time %v but returned invalid time %v", defaultTime, time)
	}
}

func TestSetPositiveIntOrZeroFromPointer(t *testing.T) {
	defaultValue := 1
	v := 10
	value := setPositiveIntOrZeroFromPointer(&v, defaultValue)
	if value == 1 {
		t.Errorf(" setPositiveIntOrZeroFromPointer() should return value %v but returned invalid value %v", value, defaultValue)
	}

	defaultValue = 1
	value = setPositiveIntOrZeroFromPointer(nil, defaultValue)
	if value != 1 {
		t.Errorf(" setPositiveIntOrZeroFromPointer() should return default value %v but returned invalid value %v", defaultValue, value)
	}
}