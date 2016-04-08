## Overview

gen-options is a small go generate tool. The project is inspired by [gen-mocks](https://github.com/sourcegraph/gen-mocks).

Suppose you have a struct looks like

```
type Foo struct {
     fstField int
     sndField string
}
```
and if you want to initialize the struct using dependency injection like this
```
func NewFoo(options ...func(*Foo)) {
     foo := &Foo{}
     for _, o := range options {
     	 o(foo)
     }
}
```
then you need to define one function per field
```
func OptionFirstField(int x) func(*Foo) {
     return func(f *Foo) {
     	    f.fstField = x
     }
}
```
gen-options can help generate these functions for you.

## Usage

install: `go get github.com/showbufire/gen-options`

help: `gen-options -h`

common usage: `gen-options -s Foo -w`

By default, for each field named `bar`, it will generate a function `OptionBar`. You can tweak the default behavior, by

* feed a `-f=Prefix` option. The generated function will be `PrefixBar`.
* add an `options` field tag to opt in for function generation.
For example, if ``bar string `options:"Baz"` ``, then the generated will be `OptionBaz`. If it's `options:""`, then the generated will be `OptionBar`.
* no function will be generated if `options` field tag is missing.
* comments will be copied from the field to the generated function with the function putting in front.
