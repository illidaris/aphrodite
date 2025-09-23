package mongoex

import (
	"context"
	"errors"
	"strings"

	"github.com/illidaris/aphrodite/component/base"
	"github.com/illidaris/aphrodite/component/embedded"
	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/illidaris/aphrodite/pkg/group"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	ErrCurNil      = errors.New("cur is nil")
	MongoComponent = embedded.NewComponent[*mongo.Client]()
	MongoNameMap   = map[string]string{}
	getKey         func(ctx context.Context) string
	clients        []*mongo.Client
)

func SetGetKeyFunc(f func(ctx context.Context) string) {
	getKey = f
}

func NewMongo(key, dbname, conn string) error {
	if !strings.Contains(conn, dbname) {
		return errors.New("dbname not in conn, please check")
	}
	opts := options.Client().ApplyURI(conn).SetLoggerOptions(NewLoggerOptions())
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return err
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return err
	}
	MongoComponent.NewWriter(key, client)
	MongoComponent.NewReader(key, client)
	MongoNameMap[key] = dbname
	clients = append(clients, client)
	return nil
}

func CloseAllMongo(ctx context.Context) {
	_, _ = group.GroupFunc(func(cs ...*mongo.Client) (int64, error) {
		for _, c := range cs {
			err := c.Disconnect(ctx)
			if err != nil {
				println("mongo client disconnect error", err.Error())
			}
		}
		return 0, nil
	}, 1, clients...)
}

// GetNamedMongoClient from mongo map
func GetNamedMongoClient(key string) *mongo.Client {
	return MongoComponent.GetWriter(key)
}

func GetMongoNameByKey(key string) string {
	c, ok := MongoNameMap[key]
	if !ok {
		return ""
	}
	return c
}

func GetMongoNameByCtx(ctx context.Context) string {
	key := getKey(ctx)
	c, ok := MongoNameMap[key]
	if !ok {
		return ""
	}
	return c
}

// SyncDbStruct
func SyncDbStruct(dbShardingKeys [][]any, pos ...dependency.IPo) error {
	return base.SyncDbStruct(func(s *base.InitTable) error {
		db := MongoComponent.GetWriter(s.Db)
		if db == nil {
			return errors.New("db is nil")
		}
		v, ok := s.P.(IRawIndex)
		if ok {
			realDBName, tbOk := MongoNameMap[s.Db]
			if !tbOk {
				return errors.New("db name not found")
			}
			_, err := db.Database(realDBName).Collection(s.Table).Indexes().CreateMany(context.Background(), v.GetRawIndexes())
			return err
		}
		return nil
	})(dbShardingKeys, pos...)
}

type IRawIndex interface {
	GetRawIndexes() []mongo.IndexModel
}
