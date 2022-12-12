package parser

import "testing"

func TestParseFiles(t *testing.T) {
	p, err := ParseSingleFile(`C:/Users/green/go/pkg/mod/github.com/!orlov!evgeny/go-mcache@v0.0.0-20200121124330-1a8195b34f3a/mcache.go`, true)
	if err != nil {
		t.Fatal(err)
	}

	s, _ := p.ImportPath()
	println(s)
}
