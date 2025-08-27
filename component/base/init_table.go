package base

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/illidaris/aphrodite/pkg/group"
	"github.com/schollz/progressbar/v3"
)

// InitTable init po table
type InitTable struct {
	Table string
	Db    string
	P     dependency.IPo
}

// SyncDbStruct
func SyncDbStruct(initFunc func(*InitTable) error) func(dbShardingKeys [][]any, pos ...dependency.IPo) error {
	return func(dbShardingKeys [][]any, pos ...dependency.IPo) error {
		ss := Trans2Table(dbShardingKeys, pos...)
		total := len(ss)
		batch := 10
		bar := progressbar.Default(int64(total))
		errMap := sync.Map{}
		affect := int32(0)
		_, _ = group.GroupFunc(func(subss ...*InitTable) (int64, error) {
			for _, s := range subss {
				err := initFunc(s)
				if err != nil {
					errMap.Store(fmt.Sprintf("[%v.%v]", s.Db, s.Table), err.Error())
				}
				_ = atomic.AddInt32(&affect, 1)
				_ = bar.Add(1)
			}
			return 0, nil
		}, batch, ss...)
		_, _ = fmt.Printf("初始化完成总计初始化[%v/%v]", affect, total)
		errMap.Range(func(key, value any) bool {
			fmt.Printf("%v初始化失败：%v\n", key, value)
			return true
		})
		return nil
	}
}

// trans2Struct transform po to table
func Trans2Table(dbShardingKeys [][]any, pos ...dependency.IPo) []*InitTable {
	ss := []*InitTable{}
	for _, p := range pos {
		dbs := ShardingDb(p, dbShardingKeys...)
		for _, db := range dbs {
			is := ShardingTable(p, db)
			ss = append(ss, is...)
		}
	}
	return ss
}

// shardingDb sharding db
func ShardingDb(p dependency.IPo, dbShardingKeys ...[]any) []string {
	dbs := []string{}
	if sharding, ok := p.(dependency.IDbSharding); ok {
		for _, dbShardingKey := range dbShardingKeys {
			db := sharding.DbSharding(dbShardingKey...)
			dbs = append(dbs, db)
		}
	} else {
		dbs = append(dbs, p.Database())
	}
	return dbs
}

// shardingTable sharding table
func ShardingTable(p dependency.IPo, db string) []*InitTable {
	ss := []*InitTable{}
	if sharding, ok := p.(dependency.ITableSharding); ok {
		for i := 0; i < int(sharding.TableTotal()); i++ {
			s := &InitTable{
				Db: db, Table: sharding.TableSharding(i), P: p,
			}
			ss = append(ss, s)
		}
	} else {
		s := &InitTable{
			Db: db, Table: p.TableName(), P: p,
		}
		ss = append(ss, s)
	}
	return ss
}
