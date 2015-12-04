package parser

import (
	"reflect"
)

type GoFile struct {
	Package string
	Structs []*GoStruct
}

type GoStruct struct {
	File   *GoFile
	Name   string
	Fields []*GoField
}

type GoField struct {
	Struct *GoStruct
	Name   string
	Type   string
	Tag    *GoTag
}

type GoTag struct {
	Field *GoField
	Value string
}

func (this *GoTag) Get(key string) string {
	return reflect.StructTag(this.Value).Get(key)
}
