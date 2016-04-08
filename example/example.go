package example

//go:generate gen-options -s Foo|Baz -w -f MyOption

import "go/ast"

type Foo struct {
	*Bar
	fst    int  `options:"First"`
	snd    *Bar `options:"Second"`
	trd    []string
	fourth *ast.Field
	fifth  int `options:"_omit"`
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
