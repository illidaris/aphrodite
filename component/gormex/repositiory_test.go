package gormex

import (
	"context"
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/IvanWhisper/aphrodite/component/dependency"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/agiledragon/gomonkey/v2"
	iLog "github.com/illidaris/logger"
	"github.com/smartystreets/goconvey/convey"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

func TestBaseRepositoryBaseCreate(t *testing.T) {
	mockDb(func(mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `test_struct`").WillReturnResult(sqlmock.NewResult(
			1, 1,
		))
		mock.ExpectCommit()
	}, func(err error) {
		if err != nil {
			t.Error(err)
		}
		ctx := context.Background()
		pos := []*testStructPo{
			{
				BizId: 1,
				Code:  "x1",
			},
		}
		convey.Convey("TestBaseRepositoryBaseCreate", t, func() {
			convey.Convey("BaseCreate", func() {
				repo := &BaseRepository[*testStructPo]{}
				affect, err := repo.BaseCreate(ctx, dependency.BaseOption{}, pos...)
				convey.So(affect, convey.ShouldEqual, 1)
				convey.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestBaseRepositoryBaseCreateIgnore(t *testing.T) {
	mockDb(func(mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT IGNORE INTO `test_struct`").WillReturnResult(sqlmock.NewResult(
			1, 1,
		))
		mock.ExpectCommit()
	}, func(err error) {
		if err != nil {
			t.Error(err)
		}
		ctx := context.Background()
		pos := []*testStructPo{
			{
				BizId: 1,
				Code:  "x1",
			},
		}
		convey.Convey("TestBaseRepositoryBaseCreate", t, func() {
			convey.Convey("BaseCreate", func() {
				repo := &BaseRepository[*testStructPo]{}
				affect, err := repo.BaseCreate(ctx, dependency.BaseOption{Ignore: true}, pos...)
				convey.So(affect, convey.ShouldEqual, 1)
				convey.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestBaseRepositoryBaseSave(t *testing.T) {
	mockDb(func(mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `test_struct`").WillReturnResult(sqlmock.NewResult(
			1, 1,
		))
		mock.ExpectCommit()
	}, func(err error) {
		if err != nil {
			t.Error(err)
		}
		ctx := context.Background()
		pos := []*testStructPo{
			{
				BizId: 1,
				Code:  "x1",
			},
		}
		convey.Convey("TestBaseRepositoryBaseSave", t, func() {
			convey.Convey("BaseSave", func() {
				repo := &BaseRepository[*testStructPo]{}
				affect, err := repo.BaseSave(ctx, dependency.BaseOption{}, pos...)
				convey.So(affect, convey.ShouldEqual, 1)
				convey.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestBaseRepositoryBaseUpdate(t *testing.T) {
	mockDb(func(mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `test_struct`").WillReturnResult(sqlmock.NewResult(
			1, 1,
		))
		mock.ExpectCommit()
	}, func(err error) {
		if err != nil {
			t.Error(err)
		}
		ctx := context.Background()
		pos := []*testStructPo{
			{
				Id:    1,
				BizId: 1,
				Code:  "x1",
			},
		}
		convey.Convey("TestBaseRepositoryBaseUpdate", t, func() {
			convey.Convey("BaseUpdate", func() {
				repo := &BaseRepository[*testStructPo]{}
				affect, err := repo.BaseUpdate(ctx, dependency.BaseOption{}, pos[0])
				convey.So(affect, convey.ShouldEqual, 1)
				convey.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestBaseRepositoryBaseSoftDelete(t *testing.T) {
	mockDb(func(mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `test_struct`").WillReturnResult(sqlmock.NewResult(
			1, 1,
		))
		mock.ExpectCommit()
	}, func(err error) {
		if err != nil {
			t.Error(err)
		}
		ctx := context.Background()
		p := &testStructDeledPo{
			Id:    1,
			BizId: 1,
			Code:  "x1",
		}
		p.Id = 1
		convey.Convey("TestBaseRepositoryBaseSoftDelete", t, func() {
			convey.Convey("BaseDelete", func() {
				repo := &BaseRepository[*testStructDeledPo]{}
				affect, err := repo.BaseDelete(ctx, dependency.BaseOption{}, p)
				convey.So(affect, convey.ShouldEqual, 1)
				convey.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestBaseRepositoryBaseDelete(t *testing.T) {
	mockDb(func(mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM `test_struct`").WillReturnResult(sqlmock.NewResult(
			1, 1,
		))
		mock.ExpectCommit()
	}, func(err error) {
		if err != nil {
			t.Error(err)
		}
		ctx := context.Background()
		pos := []*testStructPo{
			{
				Id:    1,
				BizId: 1,
				Code:  "x1",
			},
		}
		convey.Convey("TestBaseRepositoryBaseDelete", t, func() {
			convey.Convey("BaseDelete", func() {
				repo := &BaseRepository[*testStructPo]{}
				affect, err := repo.BaseDelete(ctx, dependency.BaseOption{}, pos[0])
				convey.So(affect, convey.ShouldEqual, 1)
				convey.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestBaseRepositoryBaseCount(t *testing.T) {
	mockDb(func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `test_struct`")).WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow([]driver.Value{2}...))
		mock.ExpectQuery("SELECT \\* FROM `test_struct` WHERE `test_struct`.`id` = \\?").WithArgs([]driver.Value{2}...).WillReturnRows(sqlmock.NewRows([]string{
			"id", "code",
		}).AddRow([]driver.Value{
			2, "x2",
		}...))
	}, func(err error) {
		if err != nil {
			t.Error(err)
		}
		ctx := context.Background()
		convey.Convey("TestBaseRepositoryBaseCount", t, func() {
			convey.Convey("BaseCount", func() {
				repo := &BaseRepository[*testStructPo]{}
				affect, err := repo.BaseCount(ctx, dependency.BaseOption{}, &testStructPo{})
				convey.So(affect, convey.ShouldEqual, 2)
				convey.So(err, convey.ShouldBeNil)
			})
			convey.Convey("BaseGet", func() {
				repo := &BaseRepository[*testStructPo]{}
				p := &testStructPo{Id: 2}
				affect, err := repo.BaseGet(ctx, dependency.BaseOption{}, p)
				convey.So(affect, convey.ShouldEqual, 1)
				convey.So(err, convey.ShouldBeNil)
				convey.So(p.Id, convey.ShouldEqual, 2)
			})
		})
	})
}

func TestBaseRepositoryBaseQuery(t *testing.T) {
	mockDb(func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery("SELECT \\* FROM `test_struct`").WillReturnRows(sqlmock.NewRows([]string{
			"id", "code",
		}).AddRow([]driver.Value{
			1, "x1",
		}...).AddRow([]driver.Value{
			2, "x2",
		}...).AddRow([]driver.Value{
			3, "x3",
		}...).AddRow([]driver.Value{
			4, "x4",
		}...).AddRow([]driver.Value{
			5, "x5",
		}...))

		mock.ExpectQuery("SELECT \\* FROM `test_struct` ORDER BY id LIMIT 2").WillReturnRows(sqlmock.NewRows([]string{
			"id", "code",
		}).AddRow([]driver.Value{
			2, "x2",
		}...).AddRow([]driver.Value{
			5, "x5",
		}...))

	}, func(err error) {
		if err != nil {
			t.Error(err)
		}
		ctx := context.Background()
		convey.Convey("TestBaseRepositoryBaseQuery", t, func() {
			convey.Convey("BaseQuery", func() {
				repo := &BaseRepository[*testStructPo]{}
				pos, err := repo.BaseQuery(ctx, dependency.BaseOption{}, &testStructPo{})
				convey.So(len(pos), convey.ShouldEqual, 5)
				convey.So(err, convey.ShouldBeNil)
			})
			convey.Convey("BaseQueryPage", func() {
				repo := &BaseRepository[*testStructPo]{}
				pos, err := repo.BaseQuery(ctx, dependency.BaseOption{Page: &testpage{Page: 1, Size: 2}}, &testStructPo{})
				convey.So(len(pos), convey.ShouldEqual, 2)
				convey.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestBaseRepositoryBaseQueryConds(t *testing.T) {
	mockDb(func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery("SELECT \\* FROM `test_struct` WHERE `test_struct`.`id` = \\? LIMIT 3").WithArgs([]driver.Value{2}...).WillReturnRows(sqlmock.NewRows([]string{
			"id", "code",
		}).AddRow([]driver.Value{
			1, "x2",
		}...))
		mock.ExpectQuery("SELECT `code` FROM `test_struct` WHERE `test_struct`.`id` = \\? LIMIT 1000").WithArgs([]driver.Value{2}...).WillReturnRows(sqlmock.NewRows([]string{
			"id", "code", "bizId",
		}).AddRow([]driver.Value{
			0, "x2", 0,
		}...))
		mock.ExpectQuery("SELECT `test_struct`.`id`,`test_struct`.`bizId`,`test_struct`.`status`,`test_struct`.`createBy`,`test_struct`.`createAt`,`test_struct`.`updateBy`,`test_struct`.`updateAt`,`test_struct`.`describe` FROM `test_struct` WHERE `test_struct`.`id` = \\? LIMIT 4").WithArgs([]driver.Value{2}...).WillReturnRows(sqlmock.NewRows([]string{
			"id", "code", "bizId",
		}).AddRow([]driver.Value{
			2, "", 5256,
		}...))
		mock.ExpectQuery("SELECT \\* FROM `test_struct` WHERE `test_struct`.`id` = \\? LIMIT 1000 FOR UPDATE").WithArgs([]driver.Value{2}...).WillReturnRows(sqlmock.NewRows([]string{
			"id", "code",
		}).AddRow([]driver.Value{
			3, "x3",
		}...))
	}, func(err error) {
		if err != nil {
			t.Error(err)
		}
		ctx := context.Background()
		convey.Convey("TestBaseRepositoryBaseQuery", t, func() {
			convey.Convey("BaseQueryConds", func() {
				repo := &BaseRepository[*testStructPo]{}
				pos, err := repo.BaseQuery(ctx, dependency.BaseOption{Conds: []any{2}, BatchSize: 3}, &testStructPo{})
				convey.So(len(pos), convey.ShouldEqual, 1)
				convey.So(err, convey.ShouldBeNil)
			})
			convey.Convey("BaseQuerySelectedConds", func() {
				repo := &BaseRepository[*testStructPo]{}
				pos, err := repo.BaseQuery(ctx, dependency.BaseOption{Selects: []string{`code`}, Conds: []any{2}}, &testStructPo{})
				convey.So(len(pos), convey.ShouldEqual, 1)
				convey.So(err, convey.ShouldBeNil)
				if v := pos[0]; v != nil {
					convey.So(v.Id, convey.ShouldEqual, 0)
					convey.So(v.Code, convey.ShouldEqual, "x2")
					convey.So(v.BizId, convey.ShouldEqual, 0)
				}
			})
			convey.Convey("BaseQueryOmitConds", func() {
				repo := &BaseRepository[*testStructPo]{}
				pos, err := repo.BaseQuery(ctx, dependency.BaseOption{Omits: []string{`code`}, BatchSize: 4, Conds: []any{2}}, &testStructPo{})
				convey.So(len(pos), convey.ShouldEqual, 1)
				convey.So(err, convey.ShouldBeNil)
				if v := pos[0]; v != nil {
					convey.So(v.Id, convey.ShouldEqual, 2)
					convey.So(v.Code, convey.ShouldEqual, "")
					convey.So(v.BizId, convey.ShouldEqual, 5256)
				}
			})
			convey.Convey("BaseQueryLockConds", func() {
				repo := &BaseRepository[*testStructPo]{}
				pos, err := repo.BaseQuery(ctx, dependency.BaseOption{Lock: true, Conds: []any{2}}, &testStructPo{})
				convey.So(len(pos), convey.ShouldEqual, 1)
				convey.So(err, convey.ShouldBeNil)
				if v := pos[0]; v != nil {
					convey.So(v.Id, convey.ShouldEqual, 3)
					convey.So(v.Code, convey.ShouldEqual, "x3")
				}
			})
		})
	})
}

// ================================= Mock =======================================

func mockDb(f func(sqlmock.Sqlmock), exec func(error)) {
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

	f1 := gomonkey.ApplyFunc(CoreFrmCtx, func(ctx context.Context, id string) *gorm.DB {
		return gormDb
	})
	defer f1.Reset()
	f2 := gomonkey.ApplyFunc(ReadOnly, func(ctx context.Context, id string) *gorm.DB {
		return gormDb
	})
	defer f2.Reset()
	exec(err)
}

type testStructDeledPo struct {
	Id        int64                 `json:"id" gorm:"column:id;autoIncrement;type:bigint;primaryKey;comment:唯一ID"`       // identify id
	BizId     int64                 `json:"bizId" gorm:"column:bizId;type:bigint;comment:业务"`                            // game id
	Code      string                `json:"code" gorm:"column:code;type:varchar(32);comment:编码"`                         // code
	Status    int32                 `json:"status" gorm:"column:status;type:int;default:1;comment:状态"`                   // 状态 0-默认 1-未发布 2-预发布 3-发布中 4-已结束
	CreateBy  int64                 `json:"createBy" gorm:"column:createBy;<-:create;index;type:bigint;comment:创建者"`     // 创建者
	CreateAt  int64                 `json:"createAt" gorm:"column:createAt;<-:create;index;autoCreateTime;comment:创建时间"` // 创建时间
	UpdateBy  int64                 `json:"updateBy" gorm:"column:updateBy;type:bigint;comment:修改者"`                     // 修改者
	UpdateAt  int64                 `json:"updateAt" gorm:"column:updateAt;index;autoUpdateTime;comment:修改时间"`           // 修改时间
	Describe  string                `json:"describe" gorm:"column:describe;type:varchar(255);comment:描述"`                // 描述
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:column:deleteAt"`
}

func (s testStructDeledPo) ID() any {
	return s.Id
}

func (s testStructDeledPo) TableName() string {
	return "test_struct"
}

func (s testStructDeledPo) Database() string {
	return "db"
}

type testStructPo struct {
	Id       int64  `json:"id" gorm:"column:id;autoIncrement;type:bigint;primaryKey;comment:唯一ID"`       // identify id
	BizId    int64  `json:"bizId" gorm:"column:bizId;type:bigint;comment:业务"`                            // game id
	Code     string `json:"code" gorm:"column:code;type:varchar(32);comment:编码"`                         // code
	Status   int32  `json:"status" gorm:"column:status;type:int;default:1;comment:状态"`                   // 状态 0-默认 1-未发布 2-预发布 3-发布中 4-已结束
	CreateBy int64  `json:"createBy" gorm:"column:createBy;<-:create;index;type:bigint;comment:创建者"`     // 创建者
	CreateAt int64  `json:"createAt" gorm:"column:createAt;<-:create;index;autoCreateTime;comment:创建时间"` // 创建时间
	UpdateBy int64  `json:"updateBy" gorm:"column:updateBy;type:bigint;comment:修改者"`                     // 修改者
	UpdateAt int64  `json:"updateAt" gorm:"column:updateAt;index;autoUpdateTime;comment:修改时间"`           // 修改时间
	Describe string `json:"describe" gorm:"column:describe;type:varchar(255);comment:描述"`                // 描述
}

func (s testStructPo) ID() any {
	return s.Id
}

func (s testStructPo) TableName() string {
	return "test_struct"
}

func (s testStructPo) Database() string {
	return "db"
}

type testpage struct {
	Page int
	Size int
}

func (t *testpage) GetPageIndex() int64 {
	return int64(t.Page)
}
func (t *testpage) GetPageSize() int64 {
	return int64(t.Size)
}
func (t *testpage) GetBegin() int64 {
	return (t.GetPageIndex() - 1) * t.GetPageSize()
}
func (t *testpage) GetSize() int64 {
	return int64(t.Size)
}
func (t *testpage) GetSorts() []dependency.ISortField {
	return []dependency.ISortField{&testSort{Field: "id"}}
}

type testSort struct {
	Field  string
	IsDesc bool
}

func (t *testSort) GetField() string {
	return t.Field
}
func (t *testSort) GetIsDesc() bool {
	return t.IsDesc
}
