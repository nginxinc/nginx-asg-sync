package main

import (
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

var (
	durationEscaped = strings.Join(validTimeSuffixes, "|")
	validNginxTime  = regexp.MustCompile(`^([0-9]+([` + durationEscaped + `]?){0,1} *)+$`)
)

func isValidTime(time string) bool {
	if time == "" {
		return true
	}

	time = strings.TrimSpace(time)

	return validNginxTime.MatchString(time)
}
