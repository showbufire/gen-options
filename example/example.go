package example

import "go/ast"

type Foo struct {
	fst    int
	snd    *Bar
	trd    []string
	fourth *ast.Field
}

type Bar struct {
}
