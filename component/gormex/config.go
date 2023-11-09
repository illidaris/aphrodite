package gormex

const (
	DbTXPrefix string = "_aphrodite_db_tx"
)

var (
	disableQueryFields = false
)

type ContextKey string

func SetDisableQueryFields() {
	disableQueryFields = true
}
