package typesystem

import (
	"fmt"
	"strings"
)

type ArrayValue struct {
	value []Type
}

func NewArray(value []Type) ArrayValue {
	return ArrayValue{value}
}

func (s ArrayValue) String() string {
	var b strings.Builder

	b.WriteRune('*')
	b.WriteString(fmt.Sprint(len(s.value)))
	b.WriteString("\r\n")

	for _, el := range s.value {
		b.WriteString(el.String())
	}

	return b.String()
}

func (s ArrayValue) Pretty() string {
	var b strings.Builder
	last := len(s.value) - 1

	b.WriteString("[ ")

	for i, v := range s.value {
		b.WriteString(v.Pretty())
		if i != last {
			b.WriteString(", ")
		}
	}

	b.WriteString(" ]")
	return b.String()
}

func (s ArrayValue) Value() Value {
	return s.value
}
