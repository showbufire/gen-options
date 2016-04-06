package handler

import (
	"go/parser"
	"go/token"
	"regexp"
	"testing"

	"github.com/facebookgo/ensure"
)

const (
	exampleDir    = "../example"
	exampleStruct = "Foo"
	examplePkg    = "example"
)

func TestHandlePackage(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, exampleDir, nil, parser.AllErrors)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, len(pkgs), 1)

	pat, err := regexp.Compile(exampleStruct)
	ensure.Nil(t, err)

	pkg, ok := pkgs[examplePkg]
	ensure.True(t, ok)

	tspecs := HandlePackage(pkg, pat)
	ensure.DeepEqual(t, len(tspecs), 1)
	ensure.DeepEqual(t, tspecs[0].Name.Name, exampleStruct)
}
