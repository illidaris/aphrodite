package dependency

type ISetPage interface {
	IPage
	SetPageIndex(index int64)
}

type ICondPage interface {
	IPage
	ICond
}

type ICond interface {
	IDbShardingCond
	ITbShardingCond
	GetConds() []any
}

type IDbShardingCond interface {
	GetDbShardingKeys() []any
}

type ITbShardingCond interface {
	GetTbShardingKeys() []any
}

// IPage page request
type IPage interface {
	GetPageIndex() int64
	GetPageSize() int64
	GetBegin() int64
	GetSize() int64
	GetSorts() []ISortField
}

type ISearchAfter interface {
	GetSortValues() []any
	GetSorts() []ISortField
}

type IPaginator interface {
	IPage
	GetTotalRecord() int64
	GetTotalPage() int64
}
