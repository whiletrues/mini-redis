package common

import (
	"errors"
	"testing"
)

func TestParseValidBulkString(t *testing.T) {
	// Test case 1: Valid bulk string
	buffer := []byte{'5', '\r', '\n', 'H', 'e', 'l', 'l', 'o', '\r', '\n'}

	expected1 := BulkString{Value: "Hello"}

	result1, err1 := parseBulkString(buffer)
	if err1 != nil {
		t.Errorf("Unexpected error: %v", err1)
	}
	if result1 != expected1 {
		t.Errorf("Expected %v, but got %v", expected1, result1)
	}
}

func TestParseInvalidLength(t *testing.T) {
	// Test case 2: Invalid length
	buffer := []byte{}

	_, err := parseBulkString(buffer)
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	expectedErr := errors.New("Invalid length 0")
	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error message '%v', but got '%v'", expectedErr, err)
	}
}
func TestParseSimpleError(t *testing.T) {
	// Test case 1: Valid simple error
	buffer := []byte{'E', 'R', 'R', ' ', 'S', 'o', 'm', 'e', ' ', 'e', 'r', 'r', 'o', 'r', '\r', '\n'}

	expected1 := SimpleError{Value: "ERR Some error"}

	result1, err1 := parseSimpleError(buffer)
	if err1 != nil {
		t.Errorf("Unexpected error: %v", err1)
	}
	if result1 != expected1 {
		t.Errorf("Expected %v, but got %v", expected1, result1)
	}
}

func TestParseSimpleErrorEmptyBuffer(t *testing.T) {
	// Test case 2: Empty buffer
	buffer := []byte{}

	_, err := parseSimpleError(buffer)
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	expectedErr := errors.New("Invalid length 0")
	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error message '%v', but got '%v'", expectedErr, err)
	}
}