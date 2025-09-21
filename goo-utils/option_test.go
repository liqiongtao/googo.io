package goo_utils

import (
	"fmt"
	"testing"
)

type Person struct {
	Name string
}

func NewPerson(opts ...Option[*Person]) *Person {
	u := &Person{}
	for _, opt := range opts {
		opt.Apply(u)
	}
	return u
}

func WithName(name string) Option[*Person] {
	return func(u *Person) {
		u.Name = name
	}
}

func TestNewPerson(t *testing.T) {
	u := NewPerson(WithName("liqiongtao"))
	fmt.Println(u.Name)
}
