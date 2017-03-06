package main

import (
	"reflect"
	"testing"
)

func TestDetermineUpdates(t *testing.T) {
	var tests = []struct {
		updated          []string
		nginx            []string
		expectedToAdd    []string
		expectedToDelete []string
	}{{
		updated:          []string{"10.0.0.3:80", "10.0.0.4:80"},
		nginx:            []string{"10.0.0.1:80", "10.0.0.2:80"},
		expectedToAdd:    []string{"10.0.0.3:80", "10.0.0.4:80"},
		expectedToDelete: []string{"10.0.0.1:80", "10.0.0.2:80"},
	}, {
		updated:          []string{"10.0.0.2:80", "10.0.0.3:80", "10.0.0.4:80"},
		nginx:            []string{"10.0.0.1:80", "10.0.0.2:80", "10.0.0.3:80"},
		expectedToAdd:    []string{"10.0.0.4:80"},
		expectedToDelete: []string{"10.0.0.1:80"},
	}, {
		updated: []string{"10.0.0.1:80", "10.0.0.2:80", "10.0.0.3:80"},
		nginx:   []string{"10.0.0.1:80", "10.0.0.2:80", "10.0.0.3:80"},
	}, {
	// empty values
	}}

	for _, test := range tests {
		toAdd, toDelete := determineUpdates(test.updated, test.nginx)
		if !reflect.DeepEqual(toAdd, test.expectedToAdd) || !reflect.DeepEqual(toDelete, test.expectedToDelete) {
			t.Errorf("determiteUpdates(%v, %v) = (%v, %v)", test.updated, test.nginx, toAdd, toDelete)
		}
	}
}
