package parser

import (
	"fmt"
	"go/importer"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
)

type PackImporter struct {
	Fset *token.FileSet
}

func (this *PackImporter) Import(path string) (*types.Package, error) {
	println("searching for " + path)

	pack, err := importer.Default().Import(path)

	if err != nil {
		cfg := &packages.Config{
			Fset:  this.Fset,
			Mode:  packages.NeedTypes,
			Tests: true,
		}

		// Load the package by its import path
		pkgs, err := packages.Load(cfg, path)
		if err != nil {
			return nil, err
		}

		// Check for errors
		if packages.PrintErrors(pkgs) > 0 {
			return nil, fmt.Errorf("package %s has errors", path)
		}

		// Return the first package object
		return pkgs[0].Types, nil

	}

	return pack, nil
}

//--------------------------------------------------------------------
