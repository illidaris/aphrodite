package dependency

// IPage page request
type IPage interface {
	GetPageIndex() int64
	GetPageSize() int64
	GetBegin() int64
	GetSize() int64
	GetSorts() []ISortField
}

type IPaginator interface {
	IPage
	GetTotalRecord() int64
	GetTotalPage() int64
}
