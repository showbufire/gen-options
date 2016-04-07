package handler

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"
)

func WalkPackage(pkg *ast.Package, pat *regexp.Regexp) []*ast.TypeSpec {
	var structs []*ast.TypeSpec
	if pkg.Name == "main" || strings.HasSuffix(pkg.Name, "_test") {
		return nil
	}
	ast.Walk(visitFn(func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.GenDecl:
			if node.Tok == token.TYPE {
				for _, spec := range node.Specs {
					tspec := spec.(*ast.TypeSpec)
					if _, ok := tspec.Type.(*ast.StructType); !ok {
						continue
					}
					if name := tspec.Name.Name; pat.MatchString(name) {
						structs = append(structs, tspec)
					}
				}
			}
			return false
		default:
			return true
		}
	}), pkg)
	return structs
}

type visitFn func(node ast.Node) (descend bool)

func (v visitFn) Visit(node ast.Node) ast.Visitor {
	descend := v(node)
	if descend {
		return v
	} else {
		return nil
	}
}
