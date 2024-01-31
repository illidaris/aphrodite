package structure

type ISortedUniqueArray[T any] interface {
	IUniqueArray[T]
	Swap(i, j int)
	Less(i, j int) bool
}

type IUniqueArray[T any] interface {
	Append(vs ...T)
	ToSlice() []T
	Len() int
}

// NewUniqueArray函数用于创建一个IUniqueArray[T]类型的实例。
// 该实例是一个UniqueArray[T]类型的指针，其中UniqueArray[T]是一个包含一个元素的切片和一个键为T类型值、值为空结构体的映射。
// 返回创建的实例。
func NewUniqueArray[T comparable]() IUniqueArray[T] {
	res := new(UniqueArray[T])
	res.s = []T{}
	res.m = map[T]struct{}{}
	return res
}

type UniqueArray[T comparable] struct {
	s []T
	m map[T]struct{}
}

func (n *UniqueArray[T]) Append(vs ...T) {
	for _, v := range vs {
		if _, ok := n.m[v]; !ok {
			n.s = append(n.s, v)
			n.m[v] = struct{}{}
		}
	}
}

func (n *UniqueArray[T]) ToSlice() []T {
	return n.s
}

func (n *UniqueArray[T]) Len() int {
	return len(n.s)
}

type IItemSection[T comparable] interface {
	ID() T
	Sort() int
}

// NewUniqueAnyArray创建一个新的可排序的唯一数组，数组元素类型为T，T必须是可比较的类型。
func NewUniqueAnyArray[T comparable]() ISortedUniqueArray[IItemSection[T]] {
	res := new(UniqueAnyArray[T])
	res.s = []IItemSection[T]{}
	res.m = map[T]struct{}{}
	return res
}

type UniqueAnyArray[T comparable] struct {
	s []IItemSection[T]
	m map[T]struct{}
}

func (n *UniqueAnyArray[T]) Append(vs ...IItemSection[T]) {
	for _, v := range vs {
		id := v.ID()
		if _, ok := n.m[id]; !ok {
			n.s = append(n.s, v)
			n.m[id] = struct{}{}
		}
	}
}

func (n *UniqueAnyArray[T]) ToSlice() []IItemSection[T] {
	return n.s
}

func (n *UniqueAnyArray[T]) Len() int {
	return len(n.s)
}

func (n *UniqueAnyArray[T]) Swap(i, j int) {
	n.s[i], n.s[j] = n.s[j], n.s[i]
}

func (n *UniqueAnyArray[T]) Less(i, j int) bool {
	return n.s[i].Sort() < n.s[j].Sort()
}
