package typesystem

import "fmt"

type Status struct {
	value      string
	isPositive bool
}

func NewStatus(isPositive bool, value string) Status {
	return Status{value, isPositive}
}

func (s Status) String() string {
	return fmt.Sprintf("%s\r\n", s.Value())
}

func (s Status) Pretty() string {
	prefix := "-"

	if s.isPositive {
		prefix = "+"
	}

	return fmt.Sprintf("%s %s", prefix, s.value)
}

func (s Status) Value() Value {
	prefix := "-"

	if s.isPositive {
		prefix = "+"
	}

	return fmt.Sprintf("%s%s", prefix, s.value)
}
