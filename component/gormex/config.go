package gormex

const (
	DbTXPrefix string = "_db_transaction"
)

var (
	disableQueryFields = false
)

type ContextKey string

func SetDisableQueryFields() {
	disableQueryFields = true
}
