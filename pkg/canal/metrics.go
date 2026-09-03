package canal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ReqNum = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mysql_req_num",
			Help: "The number of docs req to sync",
		}, []string{"index"},
	)
	IndexNum = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mysql_indexed_num",
			Help: "The number of docs indexed to sync",
		}, []string{"index"},
	)
	DeleteNum = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mysql_deleted_num",
			Help: "The number of docs deleted from sync",
		}, []string{"index"},
	)
)

func MetricsReqInc(key string) {
	ReqNum.WithLabelValues(key).Inc()
}

func MetricsSyncInc(key string) {
	IndexNum.WithLabelValues(key).Inc()
}

func MetricsDeleteInc(key string) {
	DeleteNum.WithLabelValues(key).Inc()
}
