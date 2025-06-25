package dto

type IRow interface {
	ToRow() []string
}

type Pager struct {
	Page
	TotalRecord int64 `json:"total"`
	TotalPage   int64 `json:"totalPage"`
}

func (r Pager) GetTotal() int64 {
	return r.TotalPage
}

func (r Pager) GetTotalRecord() int64 {
	return r.TotalRecord
}

func (r Pager) GetTotalPage() int64 {
	return r.TotalPage
}

func (p *Pager) Paginator() *Pager {
	p.TotalPage = getTotalPage(p.TotalRecord, p.GetSize())
	return p
}

// getTotalPage 计算总页数
func getTotalPage(total int64, pageSize int64) int64 {
	// 如果每页数量为0，则总页数为1
	if pageSize == 0 {
		return 1
	}
	// 如果总数量能被每页数量整除，则总页数为总数量除以每页数量
	if total%pageSize == 0 {
		return total / pageSize
	}
	// 否则，总页数为总数量除以每页数量加1
	return total/pageSize + 1
}

type RecordPager[T IRow] struct {
	Data []T `json:"data"`
	Pager
}

func (r RecordPager[T]) ToRows() [][]string {
	rows := [][]string{}
	for _, record := range r.Data {
		rows = append(rows, record.ToRow())
	}
	return rows
}

// NewRecordPager函数用于创建一个RecordPager对象。
// RecordPager对象用于分页展示记录数据。
// 参数index表示当前页码，参数size表示每页显示的记录数量，参数total表示总记录数，参数records为记录数据。
// 返回值为创建的RecordPager对象指针。
func NewRecordPager[T IRow](index, size int64, total int64, records ...T) *RecordPager[T] {
	// 创建一个RecordPager对象
	res := new(RecordPager[T])
	// 设置当前页码
	res.PageIndex = index
	// 设置每页显示的记录数量
	res.PageSize = size
	// 设置总记录数
	res.TotalRecord = total
	// 设置记录数据
	res.Data = records
	// 调用Paginator方法进行分页处理
	res.Paginator()
	// 返回创建的RecordPager对象指针
	return res
}

type RecordPtrPager[T IRow] struct {
	Data []*T `json:"data"`
	Pager
}

func NewRecordPtrPager[T IRow](index, size int64, total int64, records ...*T) *RecordPtrPager[T] {
	// 创建一个RecordPager对象
	res := new(RecordPtrPager[T])
	// 设置当前页码
	res.PageIndex = index
	// 设置每页显示的记录数量
	res.PageSize = size
	// 设置总记录数
	res.TotalRecord = total
	// 设置记录数据
	res.Data = records
	// 调用Paginator方法进行分页处理
	res.Paginator()
	// 返回创建的RecordPager对象指针
	return res
}

// Deprecate : 后期启用
type DataPager struct {
	Data interface{} `json:"data"`
	Pager
}

// Deprecate : 后期启用
// NewDataPager 返回一个 DataPager 结构体的指针，该结构体包含分页相关的信息和数据。
// 参数 data 是要分页的数据。
// 参数 index 是当前页的索引。
// 参数 size 是每页显示的记录数。
// 参数 total 是总记录数。
func NewDataPager(data interface{}, index, size, total int64) *DataPager {
	res := &DataPager{}
	res.PageIndex = index
	res.PageSize = size
	res.TotalRecord = total
	res.Data = data
	res.Paginator()
	return res
}
