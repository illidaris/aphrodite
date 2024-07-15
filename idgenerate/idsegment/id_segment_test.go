package idsegment

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/illidaris/aphrodite/idgenerate/dep"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/go-redis/redis/v8"
	"github.com/illidaris/aphrodite/component/gormex"
	"github.com/smartystreets/goconvey/convey"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestIDSegment(t *testing.T) {
	convey.Convey("TestIDSegment", t, func() {
		db, err := gorm.Open(mysql.New(mysql.Config{
			DSN:                       "root:C8PU91;FFGxE1@Pqm++;@tcp(192.168.97.71:3306)/poseidon?charset=utf8mb4&parseTime=True&loc=Local&timeout=1m", // DSN data source name
			DefaultStringSize:         256,                                                                                                              // string 类型字段的默认长度
			DisableDatetimePrecision:  true,                                                                                                             // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,                                                                                                             // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,                                                                                                             // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false,                                                                                                            // 根据当前 MySQL 版本自动配置
		}), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true, // 注意 AutoMigrate 会自动创建数据库外键约束，您可以在初始化时禁用此功能
		})
		if err != nil {
			println(err.Error())
		}
		err = db.AutoMigrate(&IdRecord{})
		if err != nil {
			println(err.Error())
		}
		repo := &IdRecordRepository{}
		dbFUnc := gomonkey.ApplyFunc(gormex.CoreFrmCtx, func(ctx context.Context, id string) *gorm.DB {
			return db
		})
		defer dbFUnc.Reset()
		cache := &IdRecordCache{}
		cache.Client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", "192.168.97.71", 6379),
			DB:       0,
			Password: "s8PPBYyhTkHBSyheLPdB",
		})
		convey.Convey("NewID", func() {
			g := &IdSegment{}
			g.Batch = 100
			g.Repo = repo
			g.Cache = cache
			var gw sync.WaitGroup
			for i := 0; i < 340; i++ {
				gw.Add(1)
				func() {
					defer gw.Done()
					seg, supseg, err := g.NewSegment(context.Background(), "testadd", dep.WithNum(17))
					if err != nil {
						t.Error(err)
					}
					if seg != nil {
						bs, _ := json.Marshal(seg)
						println(string(bs))
					}
					if supseg != nil {
						bs, _ := json.Marshal(supseg)
						println("--------------------------", string(bs))
					}
					//time.Sleep(time.Millisecond * 5)
				}()
			}
			gw.Wait()
		})
	})
}
