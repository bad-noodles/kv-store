package command

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType string

const (
	Identifier    TokenType = "Identifier"
	String        TokenType = "String"
	Integer       TokenType = "Integer"
	SquareBracket TokenType = "SquareBracket"
	Comma         TokenType = "Comma"
)

type Token struct {
	Value string
	Type  TokenType
}

type Tokenizer struct {
	input string
	index int
}

func NewTokenizer(input string) *Tokenizer {
	var tokenizer Tokenizer

	tokenizer.input = input

	return &tokenizer
}

func (t *Tokenizer) isIdentifier() bool {
	char := rune(t.input[t.index])
	return unicode.IsLetter(char) || unicode.IsDigit(char)
}

func (t *Tokenizer) swallowIdentifier() (int, int) {
	left := t.index
	for t.isIdentifier() {
		if t.index == len(t.input)-1 {
			t.index++
			break
		}

		t.index++
	}

	return left, t.index
}

func (t *Tokenizer) isWhitespace() bool {
	switch t.input[t.index] {
	case ' ', '\n':
		return true
	default:
		return false
	}
}

func (t *Tokenizer) swallowWhitespace() (int, int) {
	left := t.index
	for t.isWhitespace() {
		if t.index == len(t.input)-1 {
			t.index++
			break
		}

		t.index++
	}

	return left, t.index
}

func (t *Tokenizer) isString() bool {
	return t.input[t.index] != '"'
}

func (t *Tokenizer) isEscape() bool {
	return t.input[t.index] == '\\'
}

func (t *Tokenizer) swallowString() string {
	t.index++
	var value strings.Builder

loop:
	for {
		switch {
		case t.isEscape():
			t.index++
			value.WriteByte(t.input[t.index])
			t.index++
		case t.isString():
			value.WriteByte(t.input[t.index])
			t.index++
		default:
			break loop
		}
	}

	t.index++

	return value.String()
}

func (t *Tokenizer) isInteger() bool {
	if t.index >= len(t.input) {
		return false
	}

	return unicode.IsDigit(rune(t.input[t.index]))
}

func (t *Tokenizer) swallowInteger() string {
	left := t.index
	char := t.input[t.index]

	switch char {
	case '-', '+':
		t.index++
	}

	for t.isInteger() {
		t.index++
	}

	return t.input[left:t.index]
}

func (t *Tokenizer) NextToken() (Token, error) {
	if t.index == len(t.input) {
		return Token{}, fmt.Errorf("EOF")
	}

	if t.isWhitespace() {
		t.swallowWhitespace()
	}

	current := t.input[t.index]

	if t.isInteger() || current == '+' || current == '-' {
		return Token{
			Type:  Integer,
			Value: t.swallowInteger(),
		}, nil
	}

	if t.isIdentifier() {
		left, right := t.swallowIdentifier()

		return Token{
			Type:  Identifier,
			Value: t.input[left:right],
		}, nil
	}

	switch current {
	case '"':
		return Token{
			Type:  String,
			Value: t.swallowString(),
		}, nil
	case '[', ']':
		t.index++
		return Token{
			Type:  SquareBracket,
			Value: string(current),
		}, nil
	case ',':
		t.index++
		return Token{
			Type:  Comma,
			Value: ",",
		}, nil
	}

	return Token{}, fmt.Errorf("unexpected \"%s\" on position %d", string(t.input[t.index]), t.index)
}
