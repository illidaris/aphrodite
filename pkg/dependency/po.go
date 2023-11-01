package dependency

// IPo for db po struct
type IPo interface {
	ID() any
	TableName() string
	Database() string
}

type ITableSharding interface {
	TableSharding(keys ...any) string
}

type IDbSharding interface {
	DbSharding(keys ...any) string
}
