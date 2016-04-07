package main

import (
	"flag"
	"go/build"
	"go/parser"
	"go/token"
	"log"
	"regexp"

	"github.com/facebookgo/stackerr"
	"github.com/showbufire/gen-options/handler"
)

func main() {
	structPkgDir := flag.String("p", ".", "directory of package containing interface types")
	structPat := flag.String("s", ".+Service", "regexp pattern for selecting interface types by name")
	fset := token.NewFileSet()

	flag.Parse()

	_, err := build.Import(*structPkgDir, ".", build.FindOnly)
	if err != nil {
		log.Fatal(stackerr.Wrap(err))
	}

	pat, err := regexp.Compile(*structPat)
	if err != nil {
		log.Fatal(stackerr.Wrap(err))
	}

	pkgs, err := parser.ParseDir(fset, *structPkgDir, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(stackerr.Wrap(err))
	}

	for _, pkg := range pkgs {
		handler.WalkPackage(pkg, pat)
	}
}
