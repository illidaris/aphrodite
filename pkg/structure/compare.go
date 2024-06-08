package structure

// CompareIDArray compares the IDs in two arrays of IIDSection type, returning the number of matching IDs.
// T must be a type that supports the comparable constraint.
// src: The source array of IIDSection type.
// target: The target array of IIDSection type.
// opts: Optional functions for customizing comparison behavior.
// Returns the number of matching IDs.
func CompareIDArray[T comparable](src, target []IIDSection[T], opts ...OptionFunc[IIDSection[T]]) int64 {
	// Use CompareBaseArray for the underlying comparison logic, focusing on comparing the IDs of IIDSection.
	return CompareBaseArray[IIDSection[T]](src, target, func(s, t IIDSection[T]) bool {
		return s.ID() == t.ID()
	}, opts...)
}

// CompareArray compares the elements in two arrays, returning the number of matching elements.
// T must be a type that supports the comparable constraint.
// src: The source array.
// target: The target array.
// opts: Optional functions for customizing comparison behavior.
// Returns the number of matching elements.
func CompareArray[T comparable](src, target []T, opts ...OptionFunc[T]) int64 {
	// Use CompareBaseArray for the underlying comparison logic, focusing on comparing elements directly.
	return CompareBaseArray[T](src, target, func(s, t T) bool { return s == t }, opts...)
}

// CompareBaseArray is a generic function for comparing elements in two arrays and counting the number of matching elements.
// T: The type of elements in the array, can be any type.
// src: The source array.
// target: The target array.
// equalFunc: A function for determining if two elements are equal.
// opts: Optional functions for customizing comparison behavior.
// Returns the number of matching elements.
func CompareBaseArray[T any](src, target []T, equalFunc func(T, T) bool, opts ...OptionFunc[T]) int64 {
	// If either the source or target array is empty, there are no elements to compare, return 0 directly.
	// #01 src or target nil, than no need to compare
	if len(src) == 0 || len(target) == 0 {
		return 0
	}
	// Initialize options, processing any custom comparison options.
	// #02 get option, check params
	opt := NewOption(opts...)
	// If Max in options is 0, it means there's no limit to the number of comparisons, but in practice no comparisons are needed, return 0 directly.
	if opt.Max < 1 {
		opt.Max = int64(len(src))
	}
	// Used to count the number of matching elements.
	count := int64(0)
	// Iterate through the source array.
	for _, s := range src {
		// If the current element meets the filtering condition, skip it.
		if opt.Filtering(s) {
			continue
		}
		// If the number of already matched elements exceeds the maximum specified in the options, end the loop.
		if count > opt.Max-1 {
			break
		}
		// Iterate through the target array, looking for an element that matches the current element in the source array.
		for _, t := range target {
			// If a matching element is found, execute custom logic, increment the count, and break out of the loop.
			if equalFunc(s, t) {
				opt.Iterating(s, t)
				count++
				break
			}
		}
	}
	// Return the number of matching elements.
	return count
}
