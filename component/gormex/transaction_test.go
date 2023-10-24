package gormex

import (
	"context"
	"errors"
	"testing"

	"github.com/IvanWhisper/aphrodite/component/dependency"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/agiledragon/gomonkey/v2"
	iLog "github.com/illidaris/logger"
	"github.com/smartystreets/goconvey/convey"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestNewUnitOfWork(t *testing.T) {
	mockDbWithTrans(func(mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `test_struct`").WillReturnResult(sqlmock.NewResult(
			1, 1,
		))
		mock.ExpectExec("INSERT INTO `test_struct`").WillReturnResult(sqlmock.NewResult(
			1, 1,
		))
		mock.ExpectCommit()
	}, func(err error) {
		if err != nil {
			t.Error(err)
		}
		ctx := context.Background()

		convey.Convey("TestBaseRepositoryBaseCreate", t, func() {
			convey.Convey("BaseCreate", func() {
				uok := NewUnitOfWork("db")
				repo := &BaseRepository[*testStructPo]{}
				f1 := func(subCtx context.Context) error {
					_, err := repo.BaseUpdate(subCtx, dependency.BaseOption{}, &testStructPo{
						Id:     1,
						Code:   "122",
						Status: 2,
					})
					return err
				}
				f2 := func(subCtx context.Context) error {
					_, err := repo.BaseCreate(subCtx, dependency.BaseOption{}, &testStructPo{
						Id:   2,
						Code: "1221",
					})
					return err
				}

				err := uok.Execute(ctx, f1, f2)
				convey.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestNewUnitOfWorkRollback(t *testing.T) {
	mockDbWithTrans(func(mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `test_struct`").WillReturnResult(sqlmock.NewResult(
			1, 1,
		))
		mock.ExpectCommit()
		mock.ExpectRollback()
	}, func(err error) {
		if err != nil {
			t.Error(err)
		}
		ctx := context.Background()

		convey.Convey("TestBaseRepositoryBaseCreate", t, func() {
			convey.Convey("BaseCreate", func() {
				right := errors.New("err")
				uok := NewUnitOfWork("db")
				repo := &BaseRepository[*testStructPo]{}
				f1 := func(subCtx context.Context) error {
					_, _ = repo.BaseUpdate(subCtx, dependency.BaseOption{}, &testStructPo{
						Id:     1,
						Code:   "122",
						Status: 2,
					})
					return right
				}
				f2 := func(subCtx context.Context) error {
					_, err := repo.BaseCreate(subCtx, dependency.BaseOption{}, &testStructPo{
						Id:   2,
						Code: "1221",
					})
					return err
				}
				err := uok.Execute(ctx, f1, f2)
				convey.So(err, convey.ShouldEqual, right)
			})
		})
	})
}

func mockDbWithTrans(f func(sqlmock.Sqlmock), exec func(error)) {
	iLog.OnlyConsole()
	SetDisableQueryFields()
	db, mock, err := sqlmock.New()
	if err != nil {
		exec(err)
		return
	}
	f(mock)
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

	f3 := gomonkey.ApplyFunc(GetTransactionDb, func(id string) *gorm.DB {
		return gormDb
	})
	defer f3.Reset()

	exec(err)
}
