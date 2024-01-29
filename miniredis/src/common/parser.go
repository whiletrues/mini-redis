package common

import (
	"errors"
)

type RedisType interface {
	GetValue() interface{}
}

type SimpleString struct {
	Value string
}

func (s SimpleString) GetValue() interface{} {
	return s.Value
}

type BulkString struct {
	Value string
}

func (b BulkString) GetValue() interface{} {
	return b.Value
}

type ArrayString struct {
	Value []string
}

func (a ArrayString) GetValue() interface{} {
	return a.Value
}

type SimpleError struct {
	Value string
}

func (s SimpleError) GetValue() interface{} {
	return s.Value
}

type Integer struct {
	Value int
}

func (i Integer) GetValue() interface{} {
	return i.Value
}

type Boolean struct {
	Value bool
}

func (b Boolean) GetValue() interface{} {
	return b.Value
}

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(buffer []byte) (RedisType, error) {

	identifier := buffer[0]

	switch identifier {
	case '+':
		return parseSimpleString(buffer[1:])
	case '-':
		return parseSimpleError(buffer[1:])
	case '$':
		return parseBulkString(buffer[1:])
	case '*':
		return parseArrayString(buffer[1:])
	default:
		return nil, errors.New("Invalid identifier")
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

	countBuffer := cursor.nextLine()

	count, err := getInt(countBuffer)

	if err != nil {
		return ArrayString{}, err
	}

	values := make([]string, count)

	position := 0

	for cursor.hasNext() {

		lengthBuffer := cursor.nextLine()

		if lengthBuffer[0] != '$' {
			return ArrayString{}, errors.New("Invalid length identifier")
		}

		length, err := getInt(lengthBuffer[1:])

		if err != nil {
			return ArrayString{}, err
		}

		valueBuffer := cursor.nextLine()

		if len(valueBuffer) != length {
			return ArrayString{}, errors.New("Invalid length")
		}

		if position > len(values)-1 {
			return ArrayString{}, errors.New("Invalid length")
		}

		values[position] = string(valueBuffer)

		position++
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
