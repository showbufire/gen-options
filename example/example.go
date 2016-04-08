package example

//go:generate gen-options -s Foo|Baz -w -f MyOption

import "go/ast"

// Foo a simple struct for demo
type Foo struct {
	*Bar
	// the generated function is MyOptionFirst
	fst int `options:"First"`

	// the generated function is MyOptionSecond
	snd *Bar `options:"Second"`

	// the generated function is MyOptionTrd
	trd []string `options:""`

	// the generated function is MyOptionFourth
	// fourth field first comment 2nd line
	fourth *ast.Field `options:""`

	// there's no generated function
	fifth int
}

// Bar is another simple for demo
type Bar struct {
}

// NewFoo creates a Foo
func NewFoo(options ...func(*Foo)) *Foo {
	f := &Foo{}
	for _, o := range options {
		o(f)
	}
	return f
}
