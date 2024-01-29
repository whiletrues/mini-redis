package common

import (
	"errors"
	"reflect"
	"testing"
)

func TestParseArrayString(t *testing.T) {
	buffer := []byte("3\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$3\r\nbaz\r\n")

	expected := ArrayString{
		Value: []string{"foo", "bar", "baz"},
	}

	result, err := parseArrayString(buffer)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result. Expected: %v, Got: %v", expected, result)
	}
}

func TestParseInvalidLengthArrayString(t *testing.T) {
	buffer := []byte("2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$3\r\nbaz\r\n")

	expectedError := errors.New("Invalid length")

	_, err := parseArrayString(buffer)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(err.Error(), expectedError.Error()) {
		t.Errorf("Unexpected error. Expected: %v, Got: %v", expectedError, err)
	}
}
