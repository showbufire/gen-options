package example

//go:generate gen-options -s Foo -w

import "go/ast"

type Foo struct {
	fst    int  `options:"First"`
	snd    *Bar `options:"Second"`
	trd    []string
	fourth *ast.Field
}

type Bar struct {
}

func NewFoo(options ...func(*Foo)) *Foo {
	f := &Foo{}
	for _, o := range options {
		o(f)
	}
	return f
}
