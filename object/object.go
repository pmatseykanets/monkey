package object

import "strconv"

// Type .
type Type string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"
)

// Object .
type Object interface {
	Type() Type
	Inspect() string
}

type Integer struct {
	Value int64
}

func (*Integer) Type() Type {
	return INTEGER_OBJ
}
func (i *Integer) Inspect() string {
	return strconv.FormatInt(i.Value, 10)
}

type Boolean struct {
	Value bool
}

func (*Boolean) Type() Type {
	return BOOLEAN_OBJ
}
func (b *Boolean) Inspect() string {
	return strconv.FormatBool(b.Value)
}

type Null struct{}

func (*Null) Type() Type {
	return NULL_OBJ
}
func (*Null) Inspect() string {
	return "null"
}
