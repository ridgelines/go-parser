package parser

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type GoFile struct {
	Package         string
	Path            string
	GlobalConstants []*GoType
	GlobalVariables []*GoType
	Structs         []*GoStruct
	Interfaces      []*GoInterface
	Imports         []*GoImport
	StructMethods   []*GoStructMethod
}

func isInGoPackages(path string) bool{
	goPath := strings.Replace(os.Getenv("GOPATH"), "\\", "/", -1)
	return strings.Contains(path, goPath)
}

func (g *GoFile) ImportPath() (importPath string, isExternalPackage bool, err error) {
	isExternalPackage = false

	importPath, err = filepath.Abs(g.Path)
	if err != nil {
		return "", false, err
	}

	if _, err = os.Stat(importPath); err != nil{
		return g.Path, false, err
	}

	if !isInGoPackages(importPath){
		importPath = strings.TrimSuffix(importPath, filepath.Base(importPath))
		importPath = strings.TrimSuffix(importPath, "/")
		return importPath, false, nil
	}

	importPath, err = filepath.Abs(g.Path)
	if err != nil {
		return
	}

	importPath = strings.Replace(importPath, "\\", "/", -1)

	goPath := strings.Replace(os.Getenv("GOPATH"), "\\", "/", -1)

	isExternalPackage = true

	importPath = strings.TrimPrefix(importPath, goPath)
	importPath = strings.TrimPrefix(importPath, "/src/")
	importPath = strings.TrimPrefix(importPath, "/pkg/mod/")

	i := strings.Index(importPath, "@")
	if i > 0 {
		importPath = importPath[:i]
	}

	if strings.HasSuffix(strings.ToLower(importPath), ".go") {
		i := strings.LastIndex(importPath, "/")
		if i > 0 {
			importPath = importPath[:i]
		}
	}

	if strings.Contains(importPath, "!") { // replace "!c" to "C"
		temp := ""
		nextUppercase := false
		for i := 0; i < len(importPath); i++ {
			if importPath[i] == '!' {
				nextUppercase = true
			} else {
				if nextUppercase {
					temp += strings.ToUpper(string(importPath[i]))
					nextUppercase = false
				} else {
					temp += string(importPath[i])
				}
			}
		}
		importPath = temp
	}

	importPath = strings.TrimSuffix(importPath, "/")

	return

	//
	//importPath = strings.Replace(importPath, "\\", "/", -1)
	//
	//goPath := strings.Replace(os.Getenv("GOPATH"), "\\", "/", -1)
	//importPath = strings.TrimPrefix(importPath, goPath)
	//importPath = strings.TrimPrefix(importPath, "/src/")
	//
	//importPath = strings.TrimSuffix(importPath, filepath.Base(importPath))
	//importPath = strings.TrimSuffix(importPath, "/")
	//
	//return importPath, false, nil
}

type GoImport struct {
	File *GoFile
	Name string
	Path string
}

type GoInterface struct {
	File     *GoFile
	Name     string
	Comments string
	Methods  []*GoMethod
}

type GoMethod struct {
	Name     string
	Params   []*GoType
	Comments string
	Results  []*GoType
}

type GoStructMethod struct {
	GoMethod
	Receivers []string
}

type GoType struct {
	Name       string
	Type       string
	Underlying string
	Inner      []*GoType
}

type GoStruct struct {
	File     *GoFile
	Name     string
	Comments string
	Fields   []*GoField
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

func (g *GoTag) Get(key string) string {
	tag := strings.Replace(g.Value, "`", "", -1)
	return reflect.StructTag(tag).Get(key)
}

// For an import - guess what prefix will be used
// in type declarations.  For examples:
//    "strings" -> "strings"
//    "net/http/httptest" -> "httptest"
// Libraries where the package name does not match
// will be mis-identified.
func (g *GoImport) Prefix() string {
	if g.Name != "" {
		return g.Name
	}

	path := strings.Trim(g.Path, "\"")
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 {
		return path
	}

	return path[lastSlash+1:]
}
