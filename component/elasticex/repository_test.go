package elasticex

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/illidaris/logger"
	"github.com/olivere/elastic/v7"
)

func mockES(f func()) {
	logger.OnlyConsole()
	getClientFunc := gomonkey.ApplyFunc(CoreFrmCtx, func(ctx context.Context) *elastic.Client {
		e, err := elastic.NewClient(
			elastic.SetURL("http://192.168.97.71:11203"),
			elastic.SetErrorLog(NewLogger()),
			elastic.SetInfoLog(NewLogger()),
			elastic.SetSniff(false),
		)
		if err != nil {
			println(err.Error())
		}
		return e
	})
	defer getClientFunc.Reset()
	f()
}

func mockESX(f func()) {
	logger.OnlyConsole()
	getClientFunc := gomonkey.ApplyFunc(CoreFrmCtx, func(ctx context.Context) *elastic.Client {
		return &elastic.Client{}
	})
	defer getClientFunc.Reset()
	funcIndicesCreateService := gomonkey.ApplyMethodFunc(reflect.TypeOf(&elastic.IndicesCreateService{}), "Do", func(ctx context.Context) (*elastic.IndicesCreateResult, error) {
		return &elastic.IndicesCreateResult{
			Acknowledged:       true,
			ShardsAcknowledged: true,
			Index:              "def",
		}, nil
	})
	defer funcIndicesCreateService.Reset()
	funcIndicesExistsService := gomonkey.ApplyMethodFunc(reflect.TypeOf(&elastic.IndicesExistsService{}), "Do", func(ctx context.Context) (bool, error) {
		return true, nil
	})
	defer funcIndicesExistsService.Reset()
	funcIndexService := gomonkey.ApplyMethodFunc(reflect.TypeOf(&elastic.IndexService{}), "Do", func(ctx context.Context) (*elastic.IndexResponse, error) {
		return &elastic.IndexResponse{
			Status: 200,
		}, nil
	})
	defer funcIndexService.Reset()
	funcDeleteService := gomonkey.ApplyMethodFunc(reflect.TypeOf(&elastic.DeleteService{}), "Do", func(ctx context.Context) (*elastic.DeleteResponse, error) {
		return &elastic.DeleteResponse{
			Status: 200,
		}, nil
	})
	defer funcDeleteService.Reset()
	funcSearchService := gomonkey.ApplyMethodFunc(reflect.TypeOf(&elastic.SearchService{}), "Do", func(ctx context.Context) (*elastic.SearchResult, error) {
		bs, _ := json.Marshal(&testStructShardingPo{
			Id:   "1",
			Code: "x",
		})
		return &elastic.SearchResult{
			Hits: &elastic.SearchHits{
				TotalHits: &elastic.TotalHits{
					Value: 1,
				},
				Hits: []*elastic.SearchHit{{Source: bs}},
			},
			Status: 200,
		}, nil
	})
	defer funcSearchService.Reset()
	funcGetService := gomonkey.ApplyMethodFunc(reflect.TypeOf(&elastic.GetService{}), "Do", func(ctx context.Context) (*elastic.GetResult, error) {
		bs, _ := json.Marshal(&testStructShardingPo{
			Id:   "1",
			Code: "x",
		})
		return &elastic.GetResult{
			Found:  true,
			Source: bs,
		}, nil
	})
	defer funcGetService.Reset()
	funcCountService := gomonkey.ApplyMethodFunc(reflect.TypeOf(&elastic.CountService{}), "Do", func(ctx context.Context) (int64, error) {
		return 1, nil
	})
	defer funcCountService.Reset()
	f()
}

type testStructShardingPo struct {
	dependency.EmptyPo
	Id       string `json:"id" gorm:"column:id;autoIncrement;type:bigint;primaryKey;comment:唯一ID"`       // identify id
	BizId    int64  `json:"bizId" gorm:"column:bizId;type:bigint;comment:业务"`                            // game id
	Code     string `json:"code" gorm:"column:code;type:varchar(32);comment:编码"`                         // code
	Status   int32  `json:"status" gorm:"column:status;type:int;default:1;comment:状态"`                   // 状态 0-默认 1-未发布 2-预发布 3-发布中 4-已结束
	CreateBy int64  `json:"createBy" gorm:"column:createBy;<-:create;index;type:bigint;comment:创建者"`     // 创建者
	CreateAt int64  `json:"createAt" gorm:"column:createAt;<-:create;index;autoCreateTime;comment:创建时间"` // 创建时间
	UpdateBy int64  `json:"updateBy" gorm:"column:updateBy;type:bigint;comment:修改者"`                     // 修改者
	UpdateAt int64  `json:"updateAt" gorm:"column:updateAt;index;autoUpdateTime;comment:修改时间"`           // 修改时间
	Describe string `json:"describe" gorm:"column:describe;type:varchar(255);comment:描述"`                // 描述
}

func (s testStructShardingPo) ID() any {
	return s.Id
}

func (s testStructShardingPo) TableName() string {
	return "test_struct"
}

func (s testStructShardingPo) Database() string {
	return "db"
}

func (s testStructShardingPo) TableTotal() uint32 {
	return 20
}
func (s testStructShardingPo) TableSharding(keys ...any) string {
	if len(keys) == 0 {
		return s.TableName()
	}
	return fmt.Sprintf("%s_%v", s.TableName(), keys[0])
}
func (s testStructShardingPo) DbSharding(keys ...any) string {
	if len(keys) < 2 {
		return s.TableName()
	}
	return fmt.Sprintf("%s_%v", s.Database(), keys[1])
}
func (s testStructShardingPo) GetMapping() string {
	return `{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0
		},
		"mappings":{
			"properties":{
				"id":{
					"type":"long"
				},
				"bizId":{
					"type":"integer"
				},
				"code":{
					"type":"keyword"
				},
				"describe":{
					"type":"text"
				}
			}
		}
	}`
}
