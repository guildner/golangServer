package main

import "testing"

func TestTest(t *testing.T) {
	result := false

	if result == false {
		t.Errorf("Testing a failed test")
	}
}