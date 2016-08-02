package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
)

func ParseFile(path string) (*GoFile, error) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// File: A File node represents a Go source file: https://golang.org/pkg/go/ast/#File
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return nil, err
	}

	goFile := &GoFile{
		Path:    path,
		Package: file.Name.Name,
		Structs: []*GoStruct{},
	}

	// File.Decls: A list of the declarations in the file: https://golang.org/pkg/go/ast/#Decl
	for _, decl := range file.Decls {
		switch declType := decl.(type) {

		// GenDecl: represents an import, constant, type or variable declaration: https://golang.org/pkg/go/ast/#GenDecl
		case *ast.GenDecl:
			genDecl := declType

			// Specs: the Spec type stands for any of *ImportSpec, *ValueSpec, and *TypeSpec: https://golang.org/pkg/go/ast/#Spec
			for _, genSpec := range genDecl.Specs {
				switch genSpecType := genSpec.(type) {

				// TypeSpec: A TypeSpec node represents a type declaration: https://golang.org/pkg/go/ast/#TypeSpec
				case *ast.TypeSpec:
					typeSpec := genSpecType

					// typeSpec.Type: an Expr (expression) node: https://golang.org/pkg/go/ast/#Expr
					switch typeSpecType := typeSpec.Type.(type) {

					// StructType: A StructType node represents a struct type: https://golang.org/pkg/go/ast/#StructType
					case (*ast.StructType):
						structType := typeSpecType
						goStruct := buildGoStruct(source, goFile, typeSpec, structType)
						goFile.Structs = append(goFile.Structs, goStruct)
					// InterfaceType: An InterfaceType node represents an interface type. https://golang.org/pkg/go/ast/#InterfaceType
					case (*ast.InterfaceType):
						interfaceType := typeSpecType
						goInterface := buildGoInterface(source, goFile, typeSpec, interfaceType)
						goFile.Interfaces = append(goFile.Interfaces, goInterface)
					default:
						// a not-implemented typeSpec.Type.(type), ignore
					}
					// ImportSpec: An ImportSpec node represents a single package import. https://golang.org/pkg/go/ast/#ImportSpec
				case *ast.ImportSpec:
					importSpec := genSpec.(*ast.ImportSpec)
					goImport := buildGoImport(importSpec, goFile)
					goFile.Imports = append(goFile.Imports, goImport)
				default:
					// a not-implemented genSpec.(type), ignore
				}
			}
		default:
			// a not-implemented decl.(type), ignore
		}
	}

	return goFile, nil
}

func buildGoImport(spec *ast.ImportSpec, file *GoFile) *GoImport {
	name := ""
	if spec.Name != nil {
		name = spec.Name.Name
	}

	path := ""
	if spec.Path != nil {
		path = spec.Path.Value
	}

	return &GoImport{
		Name: name,
		Path: path,
		File: file,
	}
}

func buildGoInterface(source []byte, file *GoFile, typeSpec *ast.TypeSpec, interfaceType *ast.InterfaceType) *GoInterface {
	goInterface := &GoInterface{
		File:    file,
		Name:    typeSpec.Name.Name,
		Methods: buildMethodList(interfaceType.Methods.List, source),
	}

	return goInterface
}

func buildMethodList(fieldList []*ast.Field, source []byte) []*GoMethod {
	methods := []*GoMethod{}

	for _, field := range fieldList {
		name := getNames(field)[0]

		fType, ok := field.Type.(*ast.FuncType)
		if !ok {
			// method was not a function
			continue
		}

		goMethod := &GoMethod{
			Name:    name,
			Params:  buildTypeList(fType.Params, source),
			Results: buildTypeList(fType.Results, source),
		}

		methods = append(methods, goMethod)
	}

	return methods
}

func buildTypeList(fieldList *ast.FieldList, source []byte) []*GoType {
	types := []*GoType{}

	if fieldList != nil {
		for _, t := range fieldList.List {
			goType := buildType(t.Type, source)

			for _, n := range getNames(t) {
				copyType := copyType(goType)
				copyType.Name = n
				types = append(types, copyType)
			}
		}
	}

	return types
}

func getNames(field *ast.Field) []string {
	if field.Names == nil || len(field.Names) == 0 {
		return []string{""}
	}

	result := []string{}
	for _, name := range field.Names {
		result = append(result, name.String())
	}

	return result
}

func getTypeString(expr ast.Expr, source []byte) string {
	return string(source[expr.Pos()-1 : expr.End()-1])
}

func copyType(goType *GoType) *GoType {
	return &GoType{
		Type:  goType.Type,
		Inner: goType.Inner,
		Name:  goType.Name,
	}
}

func buildType(expr ast.Expr, source []byte) *GoType {
	innerTypes := []*GoType{}
	typeString := getTypeString(expr, source)

	switch specType := expr.(type) {
	case *ast.FuncType:
		innerTypes = append(innerTypes, buildTypeList(specType.Params, source)...)
		innerTypes = append(innerTypes, buildTypeList(specType.Results, source)...)
	case *ast.ArrayType:
		innerTypes = append(innerTypes, buildType(specType.Elt, source))
	case *ast.MapType:
		innerTypes = append(innerTypes, buildType(specType.Key, source))
		innerTypes = append(innerTypes, buildType(specType.Value, source))
	case *ast.ChanType:
		innerTypes = append(innerTypes, buildType(specType.Value, source))
	case *ast.StarExpr:
		innerTypes = append(innerTypes, buildType(specType.X, source))
	case *ast.Ellipsis:
		innerTypes = append(innerTypes, buildType(specType.Elt, source))
	case *ast.InterfaceType:
		methods := buildMethodList(specType.Methods.List, source)
		for _, m := range methods {
			innerTypes = append(innerTypes, m.Params...)
			innerTypes = append(innerTypes, m.Results...)
		}

	case *ast.Ident:
	case *ast.SelectorExpr:
	default:
		fmt.Printf("Unexpected field type: `%s`,\n %#v\n", typeString, specType)
	}

	return &GoType{
		Type:  typeString,
		Inner: innerTypes,
	}
}

func buildGoStruct(source []byte, file *GoFile, typeSpec *ast.TypeSpec, structType *ast.StructType) *GoStruct {
	goStruct := &GoStruct{
		File:   file,
		Name:   typeSpec.Name.Name,
		Fields: []*GoField{},
	}

	// Field: A Field declaration list in a struct type, a method list in an interface type,
	// or a parameter/result declaration in a signature: https://golang.org/pkg/go/ast/#Field
	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			goField := &GoField{
				Struct: goStruct,
				Name:   name.String(),
				Type:   string(source[field.Type.Pos()-1 : field.Type.End()-1]),
			}

			if field.Tag != nil {
				goTag := &GoTag{
					Field: goField,
					Value: field.Tag.Value,
				}

				goField.Tag = goTag
			}

			goStruct.Fields = append(goStruct.Fields, goField)
		}
	}

	return goStruct
}
