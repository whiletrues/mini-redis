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
	case '-':
		parseSimpleError(buffer[1:])
	case '$':
		parseBulkString(buffer[1:])
	default:
	}

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
