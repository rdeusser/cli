// package slice contains various functions for working with slices.
package slice

// ConditionFunc returns true if the condition is satisfied.
type ConditionFunc[T any] func(T) bool

// Reduce reduces a slice by ranging over it, executing fn, and if true adds it
// to a temporary slice that then gets returned.
func Reduce[T any](a []T, fn ConditionFunc[T]) []T {
	buf := make([]T, 0)

	for _, arg := range a {
		if fn(arg) {
			buf = append(buf, arg)
		}
	}

	return buf
}

// Remove cuts values out of a slice. Instead of modifying the slice, it copies
// it and returns a new one with the remaining values.
func Remove[T any](a []T, i, j int) []T {
	buf := make([]T, len(a))
	copy(buf, a)

	if i > len(buf) || j > len(buf) {
		return buf
	}

	buf = append(buf[:i], buf[j:]...)

	return buf
}

// Contains returns true if the condition is satisifed.
func Contains[T any](a []T, fn ConditionFunc[T]) bool {
	for _, arg := range a {
		if fn(arg) {
			return true
		}
	}

	return false
}
