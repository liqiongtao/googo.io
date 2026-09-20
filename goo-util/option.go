package gooutil

type Option[T any] func(t T)

func (o Option[T]) Apply(t T) {
	o(t)
}
