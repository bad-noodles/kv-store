package typesystem

import "fmt"

type BooleanValue struct {
	value bool
}

func NewBoolean(value bool) BooleanValue {
	return BooleanValue{value}
}

func (b BooleanValue) String() string {
	value := "f"
	if b.value {
		value = "t"
	}
	return fmt.Sprintf("#%s\r\n", value)
}

func (b BooleanValue) Pretty() string {
	if b.value {
		return "true"
	}

	return "false"
}

func (b BooleanValue) Value() Value {
	return b.value
}
