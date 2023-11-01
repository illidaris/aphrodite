package dependency

// IPo for db po struct
type IPo interface {
	ID() any
	TableName() string
	Database() string
}

// ITableSharding split table by keys
type ITableSharding interface {
	TableSharding(keys ...any) string
}

// IDbSharding split database by keys
type IDbSharding interface {
	DbSharding(keys ...any) string
}

// IGenerateID customer id generate
type IGenerateID interface {
	SetID(id any)
}
