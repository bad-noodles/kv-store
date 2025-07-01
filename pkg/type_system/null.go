package typesystem

type NullValue struct{}

func NewNull() NullValue {
	return NullValue{}
}

func (n NullValue) String() string {
	return "_\r\n"
}

func (n NullValue) Pretty() string {
	return "<nil>"
}

func (n NullValue) Value() Value {
	return nil
}
