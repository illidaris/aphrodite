package structure

// OptionFunc is a generic type representing functions that modify instances of the Options type.
// It accepts a pointer to an Options instance, enabling configurations on the Options instance.
type OptionFunc[T any] func(*Options[T])

// Options is a generic struct holding a set of configuration options,
// including maximum, minimum values, a filtering function, and an iterator function.
type Options[T any] struct {
	Max      int64
	Min      int64
	Filter   func(T) bool
	Iterator func(src, target T)
}

// Iterating applies the Iterator function to the given source and target objects if it's set.
// If the Iterator is not set, no operation is performed.
func (o Options[T]) Iterating(src, target T) {
	if o.Iterator != nil {
		o.Iterator(src, target)
	}
}

// Filtering checks if the given object satisfies the condition using the Filter function.
// If the Filter function is not set, it defaults to returning false.
func (o Options[T]) Filtering(src T) bool {
	if o.Filter != nil {
		return o.Filter(src)
	}
	return false
}

// NewOption creates a new Options instance and applies the provided configuration functions.
// It takes one or more OptionFunc functions as arguments to configure the Options instance.
func NewOption[T any](opts ...OptionFunc[T]) *Options[T] {
	opt := &Options[T]{}
	for _, f := range opts {
		f(opt)
	}
	return opt
}

// WithMax sets the Max property of an Options instance.
// It returns an OptionFunc function used to configure the maximum value of an Options instance.
func WithMax[T any](max int64) OptionFunc[T] {
	return func(o *Options[T]) {
		o.Max = max
	}
}

// WithMin sets the Min property of an Options instance.
// It returns an OptionFunc function used to configure the minimum value of an Options instance.
func WithMin[T any](min int64) OptionFunc[T] {
	return func(o *Options[T]) {
		o.Min = min
	}
}

// WithIterator sets the Iterator function of an Options instance.
// It returns an OptionFunc function used to configure the iteration behavior of an Options instance.
func WithIterator[T any](iterator func(s, t T)) OptionFunc[T] {
	return func(o *Options[T]) {
		o.Iterator = iterator
	}
}

// WithFilter src filter.
func WithFilter[T any](filter func(s T) bool) OptionFunc[T] {
	return func(o *Options[T]) {
		o.Filter = filter
	}
}
