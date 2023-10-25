package dependency

// ISortField sort field
type ISortField interface {
	GetField() string
	GetIsDesc() bool
}
