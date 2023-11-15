package dependency

// IPo for db po struct
type IPo interface {
	ID() any
	TableName() string
	Database() string
	ToRow() []string
	ToJson() string
}

var _ = IPo(&EmptyPo{})

// EmptyPo empty impl
type EmptyPo struct{}

func (p EmptyPo) ID() any {
	return nil
}
func (p EmptyPo) TableName() string {
	return ""
}
func (p EmptyPo) Database() string {
	return ""
}
func (p EmptyPo) ToRow() []string {
	return []string{}
}
func (p EmptyPo) ToJson() string {
	return ""
}

// ITableSharding split table by keys
type ITableSharding interface {
	TableSharding(keys ...any) string
	TableTotal() uint32
}

// IDbSharding split database by keys
type IDbSharding interface {
	DbSharding(keys ...any) string
}

// IGenerateID customer id generate
type IGenerateID interface {
	SetID(id any)
}
