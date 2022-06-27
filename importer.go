package parser

import (
	"fmt"
	"go/importer"
	"go/token"
	"go/types"
)

type PackImporter struct {
	Fset *token.FileSet
}

func (this *PackImporter) Import(path string) (*types.Package, error){
	println("searching for "+path)
	
	pack, err := importer.ForCompiler(this.Fset, "source", nil).Import(path)
	if err != nil{
		fmt.Printf("default importer: %v\n", err)
	}
	
	return pack, nil
}
//--------------------------------------------------------------------