package contextex

import "fmt"

// ContextKey is a package-specific key for values stored in a context.
type ContextKey string

// Context keys for values scoped to external services.
const (
	ElasticID ContextKey = "_aphrodite_es"
	DbTxID    ContextKey = "_aphrodite_dbtx"
	MongoID   ContextKey = "_aphrodite_mongo"
)

// ID returns a context key scoped to the supplied identifier.
func (c ContextKey) ID(id string) ContextKey {
	return ContextKey(fmt.Sprintf("%s_%s", c, id))
}
