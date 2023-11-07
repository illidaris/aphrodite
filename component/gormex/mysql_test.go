package gormex

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/smartystreets/goconvey/convey"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestNewMySqlClient(t *testing.T) {
	convey.Convey("TestNewMySqlClient", t, func() {
		convey.Convey("TestNewMySqlClient", func() {
			db, err := NewMySqlClient("", nil)
			convey.So(err, convey.ShouldBeError)
			convey.So(db.RowsAffected, convey.ShouldEqual, 0)
		})
	})
}

func TestSyncDbSruct(t *testing.T) {
	convey.Convey("TestSyncDbSruct", t, func() {
		convey.Convey("TestSyncDbSruct", func() {
			SetDisableQueryFields()
			db, _, err := sqlmock.New()
			if err != nil {
				return
			}
			defer db.Close()
			gormDb, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				DefaultStringSize:         256,  // string 类型字段的默认长度
				DisableDatetimePrecision:  true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
				DontSupportRenameIndex:    true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
				DontSupportRenameColumn:   true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
				SkipInitializeWithVersion: true, // 根据当前 MySQL 版本自动配置
			}), &gorm.Config{
				Logger:                                   NewLogger(),
				DisableForeignKeyConstraintWhenMigrating: true, // 注意 AutoMigrate 会自动创建数据库外键约束，您可以在初始化时禁用此功能
			})
			MySqlComponent.NewWriter("db", gormDb)

			dberr := SyncDbStruct(nil, &testStruct2Po{})
			convey.So(dberr, convey.ShouldBeNil)
		})
		convey.Convey("TestSyncDbSructSharding", func() {
			SetDisableQueryFields()
			db, _, err := sqlmock.New()
			if err != nil {
				return
			}
			defer db.Close()
			gormDb, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				DefaultStringSize:         256,  // string 类型字段的默认长度
				DisableDatetimePrecision:  true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
				DontSupportRenameIndex:    true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
				DontSupportRenameColumn:   true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
				SkipInitializeWithVersion: true, // 根据当前 MySQL 版本自动配置
			}), &gorm.Config{
				Logger:                                   NewLogger(),
				DisableForeignKeyConstraintWhenMigrating: true, // 注意 AutoMigrate 会自动创建数据库外键约束，您可以在初始化时禁用此功能
			})
			MySqlComponent.NewWriter("db1", gormDb)
			MySqlComponent.NewWriter("db2", gormDb)
			dberr := SyncDbStruct([][]any{{"db1"}, {"db2"}}, &testStruct2PoWithSharding{})
			convey.So(dberr, convey.ShouldBeNil)
		})
	})
}

type testStruct2Po struct {
	Id int64 `json:"id" gorm:"column:id;autoIncrement;type:bigint;primaryKey;comment:唯一ID"` // identify id
}

func (s testStruct2Po) ID() any {
	return s.Id
}

func (s testStruct2Po) TableName() string {
	return "test_struct"
}

func (s testStruct2Po) Database() string {
	return "db"
}

type testStruct2PoWithSharding struct {
	Id int64 `json:"id" gorm:"column:id;autoIncrement;type:bigint;primaryKey;comment:唯一ID"` // identify id
}

func (s testStruct2PoWithSharding) ID() any {
	return s.Id
}

func (s testStruct2PoWithSharding) TableName() string {
	return "test_struct"
}

func (s testStruct2PoWithSharding) Database() string {
	return "db"
}

func (s testStruct2PoWithSharding) TableTotal() uint32 {
	return 2
}
func (s testStruct2PoWithSharding) TableSharding(keys ...any) string {
	if len(keys) == 0 {
		return s.TableName()
	}
	return fmt.Sprintf("%s_%v", s.TableName(), keys[0])
}
func (s testStruct2PoWithSharding) DbSharding(keys ...any) string {
	if len(keys) == 0 {
		return ""
	}
	return keys[0].(string)
}
