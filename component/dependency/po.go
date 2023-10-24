package dependency

// IPo for db po struct
type IPo interface {
	ID() any
	TableName() string
	Database() string
}
