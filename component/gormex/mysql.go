package gormex

import (
	"errors"

	"github.com/illidaris/aphrodite/component/base"
	"github.com/illidaris/aphrodite/component/embedded"
	"github.com/illidaris/aphrodite/pkg/dependency"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	disableQueryFields = false
)

func SetDisableQueryFields() {
	disableQueryFields = true
}

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
	return base.SyncDbStruct(func(s *base.InitTable) error {
		db := MySqlComponent.GetWriter(s.Db)
		if db == nil {
			return errors.New("db is nil")
		}
		return db.Table(s.Table).AutoMigrate(s.P)
	})(dbShardingKeys, pos...)
}
