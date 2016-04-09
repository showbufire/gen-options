package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/tools/imports"

	"github.com/facebookgo/stackerr"
	"github.com/showbufire/gen-options/handler"
)

var (
	structPkgDir = flag.String("p", ".", "directory of package containing interface types")
	structPat    = flag.String("s", "", "regexp pattern for selecting interface types by name")
	outDir       = flag.String("o", ".", "output directory")
	writeFiles   = flag.Bool("w", false, "write over existing files in output directory (default: writes to stdout)")
	prefix       = flag.String("f", "Option", "prefix of the function names")
	typeAlias    = flag.String("t", "Option", "type alias of the return value of option functions, empty if not needed")

	fset = token.NewFileSet()
)

func main() {
	flag.Parse()

	if err := work(); err != nil {
		log.Fatal(err)
	}
}

func work() error {
	_, err := build.Import(*structPkgDir, ".", build.FindOnly)
	if err != nil {
		return stackerr.Wrap(err)
	}

	if *structPat == "" {
		return fmt.Errorf("structPat(-s option) is required")
	}
	pat, err := regexp.Compile(*structPat)
	if err != nil {
		return stackerr.Wrap(err)
	}

	pkgs, err := parser.ParseDir(fset, *structPkgDir, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return stackerr.Wrap(err)
	}

	for _, pkg := range pkgs {
		tspecs := handler.WalkPackage(pkg, pat)
		fname2gened := make(map[string][]*handler.GenResult)
		for _, tspec := range tspecs {
			gened, err := handler.GenFromStructType(*prefix, tspec, *typeAlias)
			if err != nil {
				return stackerr.Wrap(err)
			}
			filename := fset.Position(tspec.Pos()).Filename
			filename = strings.TrimSuffix(filepath.Base(filename), ".go") + "_options.go"
			if _, ok := fname2gened[filename]; !ok {
				fname2gened[filename] = []*handler.GenResult{}
			}
			for _, res := range gened {
				fname2gened[filename] = append(fname2gened[filename], res)
			}
		}
		if err := write2file(*outDir, pkg.Name, fname2gened); err != nil {
			return stackerr.Wrap(err)
		}
	}
	return nil
}

func write2file(outDir, outPkg string, fname2gened map[string][]*handler.GenResult) error {
	for filename, gened := range fname2gened {
		if len(gened) == 0 {
			continue
		}
		decls := []ast.Decl{}
		for _, g := range gened {
			decls = append(decls, g.Decl)
		}
		file := &ast.File{
			Name:  ast.NewIdent(outPkg),
			Decls: decls,
		}
		log.Println("#", filename)
		var w io.Writer
		if *writeFiles {
			if err := os.MkdirAll(outDir, 0700); err != nil {
				return stackerr.Wrap(err)
			}
			f, err := os.Create(filename)
			if err != nil {
				return stackerr.Wrap(err)
			}
			defer f.Close()
			w = f
		} else {
			w = os.Stdout
		}

		var buf bytes.Buffer
		if err := printer.Fprint(&buf, fset, file); err != nil {
			return stackerr.Wrap(err)
		}
		src := buf.Bytes()
		for _, g := range gened {
			if g.Comment == nil {
				continue
			}
			funcName := "func " + g.Name + "("
			comments := ""
			for _, c := range g.Comment.List {
				comments += c.Text + "\n"
			}
			src = bytes.Replace(src, []byte(funcName),
				[]byte(comments+funcName), 1)
		}

		src, err := imports.Process(filename, src, nil)
		if err != nil {
			return stackerr.Wrap(err)
		}

		fmt.Fprintln(w, "// generated by gen-options; DO NOT EDIT")
		fmt.Fprintln(w)
		w.Write(src)
	}
	return nil
}
