package parser

import (
	"reflect"
	"strings"
)

type GoFile struct {
	Package    string
	Structs    []*GoStruct
	Interfaces []*GoInterface
	Imports    []*GoImport
}

type GoImport struct {
	File *GoFile
	Name string
	Path string
}

type GoInterface struct {
	File    *GoFile
	Name    string
	Methods []*GoMethod
}

type GoMethod struct {
	Name    string
	Params  []*GoType
	Results []*GoType
}

type GoType struct {
	Name  string
	Type  string
	Inner []*GoType
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

// For an import - guess what prefix will be used
// in type declarations.  For examples:
//    "strings" -> "strings"
//    "net/http/httptest" -> "httptest"
// Libraries where the package name does not match
// this path will be mis-identified.
func (this *GoImport) Prefix() string {
	if this.Name != "" {
		return this.Name
	}
	path := strings.Trim(this.Path, "\"")
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 {
		return path
	}
	return path[lastSlash+1:]
}
