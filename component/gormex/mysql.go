package gormex

import (
	"errors"
	"fmt"

	"github.com/illidaris/aphrodite/component/embedded"
	"github.com/illidaris/aphrodite/pkg/dependency"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var MySqlComponent = embedded.NewComponent[*gorm.DB]()

// reference docs: https://github.com/go-sql-driver/mysql#dsn-data-source-name
// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local&timeout=1m"
// NewMySqlClient new a mysql client by dsn with logger
func NewMySqlClient(dsn string, log logger.Interface) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		Logger:                                   log,
		DisableForeignKeyConstraintWhenMigrating: true, // 注意 AutoMigrate 会自动创建数据库外键约束，您可以在初始化时禁用此功能
	})
	return db, err
}

// SyncDbStruct
func SyncDbStruct(dbShardingKeys [][]any, pos ...dependency.IPo) error {
	ss := trans2Table(dbShardingKeys, pos...)
	total := len(ss)
	for index, s := range ss {
		db := MySqlComponent.GetWriter(s.Db)
		if db == nil {
			return errors.New("db is nil")
		}
		err := db.Table(s.Table).AutoMigrate(s.P)
		_, _ = fmt.Printf("库%s表%s结构初始化[%d/%d]： 初始化，%s\n", s.Db, s.Table, index+1, total, err)
	}
	return nil
}

// initTable  init po table
type initTable struct {
	Table string
	Db    string
	P     dependency.IPo
}

// trans2Struct transform po to table
func trans2Table(dbShardingKeys [][]any, pos ...dependency.IPo) []*initTable {
	ss := []*initTable{}
	for _, p := range pos {
		dbs := shardingDb(p, dbShardingKeys...)
		for _, db := range dbs {
			is := shardingTable(p, db)
			ss = append(ss, is...)
		}
	}
	return ss
}

// shardingDb sharding db
func shardingDb(p dependency.IPo, dbShardingKeys ...[]any) []string {
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
func shardingTable(p dependency.IPo, db string) []*initTable {
	ss := []*initTable{}
	if sharding, ok := p.(dependency.ITableSharding); ok {
		for i := 1; i <= int(sharding.TableTotal()); i++ {
			s := &initTable{
				Db: db, Table: sharding.TableSharding(i), P: p,
			}
			ss = append(ss, s)
		}
	} else {
		s := &initTable{
			Db: db, Table: p.TableName(), P: p,
		}
		ss = append(ss, s)
	}
	return ss
}
