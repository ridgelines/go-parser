package parser

import (
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
		Package: file.Name.Name,
		Structs: []*GoStruct{},
	}

	// File.Decls: A list of the delcarations in the file: https://golang.org/pkg/go/ast/#Decl
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
						goFile.Structs = append(goFile.Structs, buildGoStruct(source, typeSpec, structType))
					default:
						// a not-implemented typeSpec.Type.(type), ignore
					}
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

func buildGoStruct(source []byte, typeSpec *ast.TypeSpec, structType *ast.StructType) *GoStruct {
	goStruct := &GoStruct{
		Name:   typeSpec.Name.Name,
		Fields: []*GoField{},
	}

	// Field: A Field declaration list in a struct type, a method list in an interface type,
	// or a parameter/result declaration in a signature: https://golang.org/pkg/go/ast/#Field
	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			goField := &GoField{
				Name: name.String(),
				Type: string(source[field.Type.Pos()-1 : field.Type.End()-1]),
			}

			if field.Tag != nil {
				goField.Tag = &GoTag{Value: field.Tag.Value}
			}

			goStruct.Fields = append(goStruct.Fields, goField)
		}
	}

	return goStruct
}
