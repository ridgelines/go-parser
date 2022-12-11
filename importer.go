package parser

import (
	"fmt"
	"go/importer"
	"go/token"
	"go/types"
	"io"
	"io/ioutil"
	"os"
)

type PackImporter struct {
	Fset *token.FileSet
}

func lookup(path string) (io.ReadCloser, error) {
	fullpath := "C:\\Program Files\\go\\pkg\\windows_amd64\\" + path + ".a"
	stat, err := os.Stat(fullpath)
	if os.IsNotExist(err) {
		fullpath = "c:\\src\\" + path
		stat, err = os.Stat(fullpath)
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("Did not find " + path)
		}
	}

	if stat.IsDir() {
		files, err := ioutil.ReadDir(fullpath)
		if err != nil {
			return nil, err
		}
		fullpath = fullpath + string(os.PathSeparator) + files[0].Name()

	}

	return os.Open(fullpath)
}

func (this *PackImporter) Import(path string) (*types.Package, error) {
	println("searching for " + path)

	pack, err := importer.Default().Import(path)

	if err != nil {
		pack, err = importer.ForCompiler(this.Fset, "source", nil).Import(path)
		if err != nil {
			fmt.Printf("default importer: %v\n", err)
		}
	}

	return pack, nil
}

//--------------------------------------------------------------------
