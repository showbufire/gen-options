// generated by gen-options; DO NOT EDIT

package example

import "go/ast"

func MyOptionFirst(fst int) func(*Foo) {
	return func(f *Foo) {
		f.fst = fst
	}
}

func MyOptionSecond(snd *Bar) func(*Foo) {
	return func(f *Foo) {
		f.snd = snd
	}
}

func MyOptionTrd(trd []string) func(*Foo) {
	return func(f *Foo) {
		f.trd = trd
	}
}

// MyOptionFourth fourth field first comment 1st line
// fourth field first comment 2nd line
func MyOptionFourth(fourth *ast.Field) func(*Foo) {
	return func(f *Foo) {
		f.fourth = fourth
	}
}
