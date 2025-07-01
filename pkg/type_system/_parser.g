package typesystem

import (
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
  input string
  index int
}

func NewParser(input string) *Parser {
  return &Parser{ input: input }
}

func (p *Parser) swallowUntilSep

func (p *Parser) parseStatus() (Status, error) {
  isPositive := false
  if p.input[p.index] == '+' {
    isPositive = true
  }
  p.index++


}

func (p *Parser) Parse() (Type, error) {
	switch p.input[p.index] {
	case '+':
		return NewStatus(true, p.input[0:len(input)-2]), nil
	case '-':
		return NewStatus(false, input[0:len(input)-2]), nil
	case '_':
		return NewNull(), nil
	case '$':
		data := strings.SplitN(input[1:], "\r\n", 1)
		return NewString(data[1][0 : len(data[1])-2]), nil
	case '*':
		data := strings.SplitN(input[1:], "\r\n", 1)
		length, err := strconv.Atoi(data[0])
		if err != nil {
			return nil, err
		}

    var items []Type

    for ;length > 0; length-- {
      items = append(items, ) 
    }

	default:
		return nil, fmt.Errorf("\"%s\" is not a valid data type", input)
	}
}
