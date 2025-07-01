package command

import (
	"fmt"
	"strconv"
	"strings"

	typesystem "github.com/bad-noodles/kv-store/pkg/type_system"
)

type Type string

var GET, SET Type = "GET", "SET"

type Command interface{}

type Parser struct {
	tokenizer *Tokenizer
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) parseArray() (typesystem.ArrayValue, error) {
	zero := typesystem.NewArray([]typesystem.Type{})
	var value []typesystem.Type

	for {
		item, err := p.parseValue()
		if err != nil {
			return zero, nil
		}
		value = append(value, item)

		t, err := p.assertNextToken(SquareBracket, Comma)
		if err != nil {
			return zero, err
		}

		if t.Type == SquareBracket {
			if t.Value == "]" {
				break
			}

			return zero, fmt.Errorf("expected a \"]\", but got \"%s\"", t.Value)
		}
	}

	return typesystem.NewArray(value), nil
}

func (p *Parser) assertNextToken(ts ...TokenType) (Token, error) {
	token, err := p.tokenizer.NextToken()
	if err != nil {
		return token, err
	}

	for _, t := range ts {
		if token.Type == t {
			return token, nil
		}
	}

	return token, fmt.Errorf("valid token types are \"%v\", but got \"%s\" with value \"%s\"", ts, token.Type, token.Value)
}

func (p *Parser) parseGet() (typesystem.ArrayValue, error) {
	zero := typesystem.NewArray(make([]typesystem.Type, 0))
	key, err := p.assertNextToken(Identifier)
	if err != nil {
		return zero, err
	}

	return typesystem.NewArray([]typesystem.Type{
		typesystem.NewString("GET"),
		typesystem.NewString(key.Value),
	}), nil
}

func (p *Parser) parseValue() (typesystem.Type, error) {
	token, err := p.assertNextToken(String, Integer, SquareBracket, Identifier)
	if err != nil {
		return nil, err
	}

	switch token.Type {
	case String:
		return typesystem.NewString(token.Value), nil
	case Integer:
		intValue, err := strconv.Atoi(token.Value)
		if err != nil {
			return nil, err
		}
		return typesystem.NewInteger(intValue), nil
	case SquareBracket:
		return p.parseArray()
	case Identifier:
		switch token.Value {
		case "true":
			return typesystem.NewBoolean(true), nil
		case "false":
			return typesystem.NewBoolean(false), nil
		default:
			return nil, fmt.Errorf("expected a value, got \"%s\"", token.Value)
		}
	default:
		return nil, err
	}
}

func (p *Parser) parseSet() (typesystem.ArrayValue, error) {
	zero := typesystem.NewArray(make([]typesystem.Type, 0))
	key, err := p.assertNextToken(Identifier)
	if err != nil {
		return zero, err
	}

	value, err := p.parseValue()
	if err != nil {
		return zero, err
	}

	return typesystem.NewArray([]typesystem.Type{
		typesystem.NewString("SET"),
		typesystem.NewString(key.Value),
		value,
	}), nil
}

func (p *Parser) Parse(input string) (typesystem.ArrayValue, error) {
	p.tokenizer = NewTokenizer(input)
	zero := typesystem.NewArray(make([]typesystem.Type, 0))
	identifier, err := p.assertNextToken(Identifier)
	if err != nil {
		return zero, err
	}

	switch strings.ToUpper(identifier.Value) {
	case string(GET):
		return p.parseGet()
	case string(SET):
		return p.parseSet()
	default:
		return zero, fmt.Errorf("expected a command, got \"%s\"", identifier.Value)
	}
}
