package parser

import (
	"reflect"
)

type GoFile struct {
	Package string
	Structs []*GoStruct
}

type GoStruct struct {
	Name   string
	Fields []*GoField
}

type GoField struct {
	Name string
	Type string
	Tag  *GoTag
}

type GoTag struct {
	Value string
}

func (this *GoTag) Get(key string) string {
	return reflect.StructTag(this.Value).Get(key)
}
