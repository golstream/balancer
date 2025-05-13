package utils

func Ternary[T any](
	cond bool,
	Then T,
	Else T) T {

	if cond {
		return Then
	} else {
		return Else
	}
}
