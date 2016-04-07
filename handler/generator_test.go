package handler

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"testing"

	"github.com/facebookgo/ensure"
)

func TestGenOptionFromField(t *testing.T) {
	funcDecl := genOptionFromField("Foo", &ast.Field{
		Type:  ast.NewIdent("Bar"),
		Names: []*ast.Ident{ast.NewIdent("b")},
	})
	var buf bytes.Buffer
	err := printer.Fprint(&buf, token.NewFileSet(), funcDecl)
	ensure.Nil(t, err)
	fmt.Println(buf.String())
}
