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

type Member struct {
	Name string
}

func NewMember(opts ...Option[*Member]) *Member {
	u := &Member{}
	for _, opt := range opts {
		opt.Apply(u)
	}
	return u
}

func WithName[T any](name string) Option[T] {
	return func(u T) {
		if p, ok := any(u).(*Person); ok {
			p.Name = name
		}
		if m, ok := any(u).(*Member); ok {
			m.Name = name
		}
	}
}

func TestNewPerson(t *testing.T) {
	u := NewPerson(WithName[*Person]("liqiongtao"))
	fmt.Println(u.Name)

	m := NewMember(WithName[*Member]("liqiongtao"))
	fmt.Println(m.Name)
}
