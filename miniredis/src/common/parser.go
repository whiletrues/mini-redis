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
	case '$':
		parseBulkString(buffer[1:])
	default:
	}

}

type Cursor struct {
	buffer []byte
	Index  int
}

func (c *Cursor) Next() {
	c.Index++
}

func (c *Cursor) Previous() {
	c.Index--
}

func (c *Cursor) Current() (int, byte) {
	return c.Index, c.buffer[c.Index]
}

func (c *Cursor) hasNext() bool {
	return c.Index < len(c.buffer)
}

func (c *Cursor) nextLine() []byte {

	start := c.Index

	for c.hasNext() {
		index, char := c.Current()

		if char == '\r' {
			c.Next()
			c.Next()
			return c.buffer[start:index]
		}
		c.Next()
	}

	return c.buffer[start:c.Index]
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
