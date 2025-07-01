package typesystem

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type Parser struct {
	reader io.Reader
	data   Type
	error  error
	sep    []byte
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{reader: reader, sep: []byte("\r\n")}
}

func (c *Parser) swallowSeparator() error {
	data := make([]byte, 2)
	_, err := c.reader.Read(data)
	if err != nil {
		return err
	}

	if !bytes.Equal(data, c.sep) {
		return fmt.Errorf("expected \"\\r\\n\" but got \"%s\"", string(data))
	}

	return nil
}

func (c *Parser) swallowTypeId() (rune, error) {
	data := make([]byte, 1)
	_, err := c.reader.Read(data)
	if err != nil {
		return ' ', err
	}

	return rune(data[0]), nil
}

func (c *Parser) swallowUntilSeparator() (string, error) {
	var swallowed []byte
	var size int

	for {
		data := make([]byte, 1)
		_, err := c.reader.Read(data)
		if err != nil {
			return "", err
		}

		size++
		swallowed = append(swallowed, data[0])

		if size > 2 && bytes.Equal(swallowed[size-2:], c.sep) {
			break
		}
	}

	return string(swallowed[:size-2]), nil
}

func (c *Parser) swallowString() (string, error) {
	sizeStr, err := c.swallowUntilSeparator()
	if err != nil {
		return "", err
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return "", err
	}

	data := make([]byte, size)
	_, err = c.reader.Read(data)
	if err != nil {
		return "", err
	}

	err = c.swallowSeparator()

	return string(data), err
}

func (c *Parser) swallowArray() ([]Type, error) {
	zero := make([]Type, 0)
	sizeStr, err := c.swallowUntilSeparator()
	if err != nil {
		return zero, err
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return zero, err
	}

	items := make([]Type, 0, size)

	for ; size > 0; size-- {
		item, err := c.swallowNext()
		if err != nil {
			return items, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (p *Parser) swallowBoolean() (bool, error) {
	value, err := p.swallowUntilSeparator()
	if err != nil {
		return false, err
	}

	switch value {
	case "t":
		return true, nil
	case "f":
		return false, nil
	default:
		return false, fmt.Errorf("\"%s\" is not a valid value for boolean value", value)
	}
}

func (p *Parser) Next() bool {
	p.data, p.error = p.swallowNext()

	if p.error != nil && p.error.Error() == "EOF" {
		p.error = nil
	}

	return p.data != nil && p.error == nil
}

func (c *Parser) Data() Type {
	return c.data
}

func (c *Parser) Error() error {
	return c.error
}

func (c *Parser) swallowNext() (Type, error) {
	typeId, err := c.swallowTypeId()
	if err != nil {
		return nil, err
	}

	switch typeId {
	case '+':
		status, err := c.swallowUntilSeparator()
		if err != nil {
			return nil, err
		}

		return NewStatus(true, status), nil
	case '-':
		status, err := c.swallowUntilSeparator()
		if err != nil {
			return nil, err
		}

		return NewStatus(false, status), nil
	case '_':
		err := c.swallowSeparator()
		if err != nil {
			return nil, err
		}

		return NewNull(), nil
	case '$':
		value, err := c.swallowString()
		if err != nil {
			return nil, err
		}

		return NewString(value), nil
	case '*':
		value, err := c.swallowArray()
		if err != nil {
			return nil, err
		}

		return NewArray(value), nil
	case '#':
		value, err := c.swallowBoolean()
		if err != nil {
			return nil, err
		}

		return NewBoolean(value), nil
	}

	return nil, fmt.Errorf("\"%s\" is not a valid data type identifier", string(typeId))
}
