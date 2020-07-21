package schema

import "strings"

var Types = struct {
	Int32  string
	Int64  string
	UInt32 string
	UInt64 string
	String string
}{
	Int32:  "int32",
	Int64:  "int", // todo сделать явно int64, надо фиксить в lomik'e
	UInt32: "uint32",
	UInt64: "uint",  // todo сделать явно uint64, надо фиксить в lomik'e
	String: "string",
}

type Field struct {
	Name       string
	Type       string
	FieldNo    int
	PackFunc   string
	UnpackFunc string
	Size       int
}

type Index struct {
	Name         string
	Fields       []string
	FieldsStucts []Field
	Uniq         bool
	Type         string
	IndexNo      int
}

type BoxSchema struct {
	BoxName         string
	Package         string
	Space           int
	Fields          []Field
	Indexes         []Index
	PrimaryIndex    Index
	SecondaryFields []Field
}

func validate(name string, space int, fields []Field, indexes []Index) (BoxSchema, error) {
	// todo validate fields and indexes has uniq names
	// todo correct setters fields
	return BoxSchema{name, strings.ToLower(name), space, fields, indexes, indexes[0], fields[1:]}, nil
}

func (bs BoxSchema) BoxConfig() string {
	// todo box.cfg
	return ""
}
