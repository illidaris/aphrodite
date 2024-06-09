package structure

// IIDSection[T comparable] id section
type IIDSection[T comparable] interface {
	ID() T
}

func NewUnqueFilter[T comparable]() func(T) bool {
	uniqueMap := map[T]struct{}{}
	return func(id T) bool {
		_, ok := uniqueMap[id]
		if !ok {
			uniqueMap[id] = struct{}{}
		}
		return ok
	}
}
