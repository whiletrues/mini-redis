package common

import (
	"errors"
	"reflect"
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

func TestParseArrayString(t *testing.T) {

	buffer := []byte{'3', '\r', '\n', 'H', 'e', 'l', 'l', 'o', '\r', '\n', '!', '\r', '\n'}

	expected1 := ArrayString{Value: []Value{"Hello", "World", "!"}}

	result1, err1 := parseArrayString(buffer)
	if err1 != nil {
		t.Errorf("Unexpected error: %v", err1)
	}
	if reflect.DeepEqual(result1.Value, expected1.Value) {
		t.Errorf("Expected %v, but got %v", expected1, result1)
	}
}
func TestParseSimpleString(t *testing.T) {
	buffer := []byte{'H', 'e', 'l', 'l', 'o', '\r', '\n'}

	expected1 := SimpleString{Value: "Hello"}

	result1, err1 := parseSimpleString(buffer)
	if err1 != nil {
		t.Errorf("Unexpected error: %v", err1)
	}
	if result1 != expected1 {
		t.Errorf("Expected %v, but got %v", expected1, result1)
	}
}

func TestParseEmptySimpleString(t *testing.T) {

	buffer := []byte{}

	_, err2 := parseSimpleString(buffer)
	if err2 == nil {
		t.Error("Expected error, but got nil")
	}
	expectedErr2 := errors.New("Invalid length 0")
	if err2.Error() != expectedErr2.Error() {
		t.Errorf("Expected error message '%v', but got '%v'", expectedErr2, err2)
	}
}
