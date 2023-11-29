package elasticex

import "fmt"

const (
	ESPrefix string = "_aphrodite_es"
)

type ContextKey string

func GetDbTX(id string) ContextKey {
	return ContextKey(fmt.Sprintf("%s_%s", ESPrefix, id))
}
