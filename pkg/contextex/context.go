package contextex

import "fmt"

type ContextKey string

const (
	ElasticID ContextKey = "_aphrodite_es"
	DbTxID    ContextKey = "_aphrodite_dbtx"
)

func (c ContextKey) ID(id string) ContextKey {
	return ContextKey(fmt.Sprintf("%s_%s", c, id))
}
