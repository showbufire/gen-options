package example

type Bar struct {
}

type Foo struct {
	fst int
	snd *Bar
}

type Option func(*Foo)

func OptionFst(fst int) Option {
	return func(f *Foo) {
		f.fst = fst
	}
}

func OptionSnd(snd *Bar) Option {
	return func(f *Foo) {
		f.snd = snd
	}
}

func NewFoo(options ...Option) *Foo {
	f := &Foo{}
	for _, o := range options {
		o(f)
	}
	return f
}
