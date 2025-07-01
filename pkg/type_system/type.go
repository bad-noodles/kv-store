package typesystem

type Value interface {
	any | string | []Type
}

type Type interface {
	String() string
	Pretty() string
	Value() Value
}
