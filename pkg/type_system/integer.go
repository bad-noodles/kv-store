package typesystem

import "fmt"

type IntegerValue struct {
	value int
}

func NewInteger(value int) IntegerValue {
	return IntegerValue{value}
}

func (i IntegerValue) String() string {
	return fmt.Sprintf(":%d\r\n", i.value)
}

func (i IntegerValue) Pretty() string {
	return fmt.Sprintf("%d", i.value)
}

func (i IntegerValue) Value() Value {
	return i.value
}
