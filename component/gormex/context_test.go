package gormex

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/smartystreets/goconvey/convey"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestNewContext(t *testing.T) {
	convey.Convey("TestNewContext", t, func() {
		convey.Convey("TestNewContextNil", func() {
			key := "def1"
			ctx := context.Background()
			db := &gorm.DB{RowsAffected: 1}
			MySqlComponent.NewWriter(key, db)
			ctx = NewContext(ctx, key, nil)
			v := ctx.Value(GetDbTX(key))
			d := v.(*gorm.DB)
			convey.So(d.RowsAffected, convey.ShouldEqual, db.RowsAffected)
		})
		convey.Convey("TestNewContextNew", func() {
			key := "def2"
			ctx := context.Background()
			db := &gorm.DB{RowsAffected: 1}
			ctx = NewContext(ctx, key, db)
			v := ctx.Value(GetDbTX(key))
			d := v.(*gorm.DB)
			convey.So(d.RowsAffected, convey.ShouldEqual, db.RowsAffected)
		})
	})
}

func TestWithContext(t *testing.T) {
	convey.Convey("TestWithContext", t, func() {
		convey.Convey("TestNewContextNil", func() {
			key := "def1"
			ctx := context.Background()
			db := &gorm.DB{RowsAffected: 1}
			MySqlComponent.NewWriter(key, db)
			ctx = NewContext(ctx, key, nil)
			d := WithContext(ctx, key)
			convey.So(d.RowsAffected, convey.ShouldEqual, db.RowsAffected)
		})
		convey.Convey("TestNewContextNew", func() {
			key := "def2"
			ctx := context.Background()
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
			MySqlComponent.NewWriter(key, gormDb)
			d := WithContext(ctx, key)
			convey.So(d.Config.Name(), convey.ShouldEqual, gormDb.Config.Name())
		})
	})
}
