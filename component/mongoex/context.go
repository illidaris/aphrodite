package mongoex

import "github.com/illidaris/aphrodite/pkg/contextex"

func GetDbTX(id string) contextex.ContextKey {
	return contextex.DbTxID.ID(id)
}
