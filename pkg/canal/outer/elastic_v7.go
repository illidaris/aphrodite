package outer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/illidaris/aphrodite/pkg/canal"
)

var clientsSm sync.RWMutex // 嵌入读写锁
var clients = map[string]*elasticsearch.Client{}

type ElasticOuterOptionFunc func(*ElasticOuterOption)

func NewElasticOuterOption(opts ...ElasticOuterOptionFunc) *ElasticOuterOption {
	option := &ElasticOuterOption{
		EsUrls:        []string{""},
		Index:         "",
		NumWorkers:    4,
		FlushBytes:    5 * 1024 * 1024,
		FlushInterval: 1 * time.Second,
		BaseOuter: canal.BaseOuter{
			Log: canal.DefaultLogger{},
		},
	}
	return option.WithOption(opts...)
}

type ElasticOuterOption struct {
	Name            string        `json:"name" yaml:"name"`
	EsUrls          []string      `json:"es_urls" yaml:"es_urls"`
	ESUser          string        `json:"es_user" yaml:"es_user"`
	ESPwd           string        `json:"es_pwd" yaml:"es_pwd"`
	Index           string        `json:"index" yaml:"index"`
	NumWorkers      int           `json:"num_workers" yaml:"num_workers"`       // 并发工作协程数量。
	FlushBytes      int           `json:"flush_bytes" yaml:"flush_bytes"`       // 按数据体积触发刷新的阈值（字节）。
	FlushInterval   time.Duration `json:"flush_interval" yaml:"flush_interval"` // 按时间触发刷新的间隔。
	canal.BaseOuter `json:"-" yaml:"-"`
}

func (i *ElasticOuterOption) WithOption(opts ...ElasticOuterOptionFunc) *ElasticOuterOption {
	for _, opt := range opts {
		if opt != nil {
			opt(i)
		}
	}
	return i
}

func WithName(v string) ElasticOuterOptionFunc { return func(o *ElasticOuterOption) { o.Name = v } }
func WithEsUrls(v ...string) ElasticOuterOptionFunc {
	return func(o *ElasticOuterOption) { o.EsUrls = v }
}
func WithESUser(v string) ElasticOuterOptionFunc { return func(o *ElasticOuterOption) { o.ESUser = v } }
func WithESPwd(v string) ElasticOuterOptionFunc  { return func(o *ElasticOuterOption) { o.ESPwd = v } }
func WithIndex(v string) ElasticOuterOptionFunc  { return func(o *ElasticOuterOption) { o.Index = v } }
func WithNumWorkers(v int) ElasticOuterOptionFunc {
	return func(o *ElasticOuterOption) { o.NumWorkers = v }
}
func WithFlushBytes(v int) ElasticOuterOptionFunc {
	return func(o *ElasticOuterOption) { o.FlushBytes = v }
}
func WithFlushInterval(v time.Duration) ElasticOuterOptionFunc {
	return func(o *ElasticOuterOption) { o.FlushInterval = v }
}
func WithBaseOuter(v canal.BaseOuter) ElasticOuterOptionFunc {
	return func(o *ElasticOuterOption) { o.BaseOuter = v }
}

func (i ElasticOuterOption) EsCfg() elasticsearch.Config {
	return elasticsearch.Config{
		Addresses: i.EsUrls,
		Username:  i.ESUser,
		Password:  i.ESPwd,
	}
}

func (i ElasticOuterOption) BulkCfg(es *elasticsearch.Client) esutil.BulkIndexerConfig {
	return esutil.BulkIndexerConfig{
		Client:        es,
		Index:         i.Index,
		NumWorkers:    i.NumWorkers,
		FlushBytes:    i.FlushBytes,
		FlushInterval: i.FlushInterval,
	}
}

func NewElastic(option *ElasticOuterOption) (*elasticsearch.Client, error) {
	clientsSm.Lock()
	defer clientsSm.Unlock()
	c, ok := clients[option.Name]
	if ok {
		return c, nil
	}
	esClient, err := elasticsearch.NewClient(option.EsCfg())
	if err != nil {
		return nil, err
	}
	clients[option.Name] = esClient
	return esClient, nil
}

func NewElasticOuter(opts ...ElasticOuterOptionFunc) (*ElasticOuter, error) {
	option := NewElasticOuterOption(opts...)
	outer := &ElasticOuter{}
	esClient, err := NewElastic(option)
	if err != nil {
		return nil, err
	}
	bIdx, err := esutil.NewBulkIndexer(option.BulkCfg(esClient))
	if err != nil {
		return nil, err
	}
	outer.Core = bIdx
	return outer, nil
}

var _ = canal.IOuter(ElasticOuter{})

type ElasticOuter struct {
	Core esutil.BulkIndexer
}

func (i ElasticOuter) Stats() canal.OperateStats {
	s := i.Core.Stats()
	return convertStatsToOperateStats(s)
}

func (i ElasticOuter) Close(ctx context.Context) error {
	return i.Core.Close(ctx)
}

func (i ElasticOuter) GetEs(ctx context.Context, key string) (*elasticsearch.Client, error) {
	es, ok := clients[key]
	if !ok {
		return es, errors.New("no found es client")
	}
	return es, nil
}

func (i ElasticOuter) Check(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	es, err := i.GetEs(ctx, key)
	if err != nil {
		return err
	}
	res, err := es.Ping(es.Ping.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("%s状态%s", key, res.Status())
	}
	return nil
}

func (i ElasticOuter) SyncStruct(ctx context.Context, key, index, mapping string) error {
	es, err := i.GetEs(ctx, key)
	if err != nil {
		return err
	}
	createRes, err := es.Indices.Create(
		index,
		es.Indices.Create.WithContext(ctx),
		es.Indices.Create.WithBody(strings.NewReader(mapping)), // 核心步骤
	)
	if err != nil {
		return err
	}
	if createRes.IsError() {
		return fmt.Errorf("错误原因: %s", createRes.String())
	}
	return nil
}

func (i ElasticOuter) Sync(ctx context.Context, entries ...canal.Entry) (bool, error) {
	var (
		wg     sync.WaitGroup
		failed atomic.Bool
	)
	for _, entry := range entries {
		item := esutil.BulkIndexerItem{
			Index:      entry.Index,
			Action:     string(entry.Act),
			DocumentID: entry.Id,
			OnSuccess: func(_ context.Context, indexItem esutil.BulkIndexerItem, _ esutil.BulkIndexerResponseItem) {
				if indexItem.Action != string(canal.ActionDelete) {
					canal.MetricsSyncInc(indexItem.Index)
				} else {
					canal.MetricsDeleteInc(indexItem.Index)
				}
				wg.Done()
			},
			OnFailure: func(_ context.Context, _ esutil.BulkIndexerItem, _ esutil.BulkIndexerResponseItem, _ error) {
				failed.Store(true)
				wg.Done()
			},
		}
		// 非删除操作才需要设置Body
		if entry.Act != canal.ActionDelete {
			payload, err := json.Marshal(entry.Doc)
			if err != nil {
				return false, fmt.Errorf("doc marshal %s", err)
			}
			item.Body = bytes.NewReader(payload)
		}
		canal.MetricsReqInc(entry.Index)
		wg.Add(1)
		if err := i.Core.Add(ctx, item); err != nil {
			wg.Done()
			return false, err
		}
	}
	wg.Wait()
	return !failed.Load(), nil
}

func convertStatsToOperateStats(stats esutil.BulkIndexerStats) canal.OperateStats {
	return canal.OperateStats{
		NumAdded:    stats.NumAdded,
		NumFlushed:  stats.NumFlushed,
		NumFailed:   stats.NumFailed,
		NumIndexed:  stats.NumIndexed,
		NumCreated:  stats.NumCreated,
		NumUpdated:  stats.NumUpdated,
		NumDeleted:  stats.NumDeleted,
		NumRequests: stats.NumRequests,
	}
}
