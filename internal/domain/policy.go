package domain

type Policy[T any] interface {
	Check(credentials Principal[T]) error
}
