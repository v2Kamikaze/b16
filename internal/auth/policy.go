package auth

type Policy[T any] interface {
	Check(principal Principal[T]) error
}
