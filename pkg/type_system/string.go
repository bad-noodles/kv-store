package typesystem

import "fmt"

type StringValue struct {
	value string
	size  int
}

func NewString(value string) StringValue {
	return StringValue{value, len([]byte(value))}
}

func (s StringValue) String() string {
	return fmt.Sprintf("$%d\r\n%s\r\n", s.size, s.value)
}

func (s StringValue) Pretty() string {
	return fmt.Sprintf("\"%s\"", s.value)
}

func (s StringValue) Value() Value {
	return s.value
}
