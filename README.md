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
* add an `options` field tag to change the suffix of the function name.
For example, if `Bar string ``options:"Baz"```, then the generated will be `OptionBaz`.
* omit a field by adding omit tag ```options:"_omit"```. No function will be generated for this field.
