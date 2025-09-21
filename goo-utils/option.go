package goo_utils

type Option[T any] func(t T)

func (o Option[T]) Apply(t T) {
	o(t)
}
