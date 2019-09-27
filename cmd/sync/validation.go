package main

import (
	"errors"
	"regexp"
	"strings"
)

// http://nginx.org/en/docs/syntax.html
var validTimeSuffixes = []string{
	"ms",
	"s",
	"m",
	"h",
	"d",
	"w",
	"M",
	"y",
}

var durationEscaped = strings.Join(validTimeSuffixes, "|")
var validNginxTime = regexp.MustCompile(`^([0-9]+([` + durationEscaped + `]?){0,1} *)+$`)

func validateTime(time string) error {
	if time == "" {
		return nil
	}

	if _, err := ParseTime(time); err != nil {
		return err
	}

	return nil
}

// ParseTime ensures that the string value in the annotation is a valid time.
func ParseTime(s string) (string, error) {
	s = strings.TrimSpace(s)

	if validNginxTime.MatchString(s) {
		return s, nil
	}
	return "", errors.New("Invalid time string")
}

