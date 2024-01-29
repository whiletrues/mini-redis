package common

import (
	"errors"
)

type Value interface{}

type SimpleString struct {
	Value string
}

type BulkString struct {
	Value string
}

type ArrayString struct {
	Value []Value
}

type SimpleError struct {
	Value string
}

type Integer struct {
	Value int
}

type Boolean struct {
	Value bool
}

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(buffer []byte) {

	identifier := buffer[0]

	switch identifier {
	case '+':
		parseSimpleString(buffer[1:])
	case '-':
		parseSimpleError(buffer[1:])
	case '$':
		parseBulkString(buffer[1:])
	default:
	}
}

func parseSimpleString(buffer []byte) (SimpleString, error) {

	cursor := Cursor{
		buffer: buffer,
		Index:  0,
	}

	valueBuffer := cursor.nextLine()

	if len(valueBuffer) == 0 {
		return SimpleString{}, errors.New("Invalid length 0")
	}

	return SimpleString{Value: string(valueBuffer)}, nil
}

func parseArrayString(buffer []byte) (ArrayString, error) {

	cursor := Cursor{
		buffer: buffer,
		Index:  0,
	}

	lengthBuffer := cursor.nextLine()

	length, err := getInt(lengthBuffer)

	if err != nil {
		panic(err)
	}
	values := make([]Value, length)

	for cursor.hasNext() {

		valueBuffer := cursor.nextLine()

		value := string(valueBuffer)

		values = append(values, value)
	}

	return ArrayString{Value: values}, nil
}

func getInt(buffer []byte) (int, error) {
	var x int32
	for _, c := range buffer {
		x = x*10 + int32(c-'0')
	}
	return int(x), nil
}

func parseSimpleError(buffer []byte) (SimpleError, error) {

	cursor := Cursor{
		buffer: buffer,
		Index:  0,
	}

	valueBuffer := cursor.nextLine()

	return SimpleError{Value: string(valueBuffer)}, nil
}

func parseBulkString(buffer []byte) (BulkString, error) {

	cursor := Cursor{
		buffer: buffer,
		Index:  0,
	}

	lengthBuffer := cursor.nextLine()

	//length := int(lengthBuffer)

	if len(lengthBuffer) == 0 {
		return BulkString{}, errors.New("Invalid length 0")
	}

	valueBuffer := cursor.nextLine()

	return BulkString{Value: string(valueBuffer)}, nil
}

func isNumeric(char byte) bool {
	return char >= '0' && char <= '9'
}

func isAlphaNumeric(char byte) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9'
}

func isAlpha(char byte) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z'
}
