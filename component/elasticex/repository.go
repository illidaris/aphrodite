package elasticex

import (
	"github.com/illidaris/aphrodite/pkg/dependency"
)

var _ = dependency.IRepository[dependency.IEntity](&BaseRepository[dependency.IEntity]{}) // impl check

type BaseRepository[T dependency.IEntity] struct{} // base repository
