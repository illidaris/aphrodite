package elasticex

import (
	"github.com/illidaris/aphrodite/component/embedded"
	"github.com/olivere/elastic/v7"
)

var ElasticComponent = embedded.NewComponent[*elastic.Client]()

func NewElasticClient(urls ...string) (*elastic.Client, error) {
	return elastic.NewClient(
		elastic.SetURL(urls...),
		elastic.SetErrorLog(NewLogger()),
		elastic.SetInfoLog(NewLogger()))
}
