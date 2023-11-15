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

func (p *Pager) Paginator() *Pager {
	p.TotalPage = getTotalPage(p.TotalRecord, p.GetSize())
	return p
}

func getTotalPage(total int64, pageSize int64) int64 {
	if pageSize == 0 {
		return 1
	}
	if total%pageSize == 0 {
		return total / pageSize
	}
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

func NewRecordPager[T IRow](index, size int64, total int64, records ...T) *RecordPager[T] {
	res := new(RecordPager[T])
	res.PageIndex = index
	res.PageSize = size
	res.TotalRecord = total
	res.Data = records
	res.Paginator()
	return res
}

// Deprecate : 后期启用
type DataPager struct {
	Data interface{} `json:"data"`
	Pager
}

// Deprecate : 后期启用
func NewDataPager(data interface{}, index, size int64, total int64) *DataPager {
	res := &DataPager{}
	res.PageIndex = index
	res.PageSize = size
	res.TotalRecord = total
	res.Data = data
	res.Paginator()
	return res
}
