package mongoex

import (
	"context"
	"errors"

	"github.com/illidaris/aphrodite/component/embedded"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	ErrCurNil      = errors.New("cur is nil")
	MongoComponent = embedded.NewComponent[*mongo.Client]()
	MongoNameMap   = map[string]string{}
	getKey         func(ctx context.Context) string
)

func SetGetKeyFunc(f func(ctx context.Context) string) {
	getKey = f
}

func NewMongo(key, dbname, conn string) {
	opts := options.Client().ApplyURI(conn).SetLoggerOptions(NewLoggerOptions())
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		l.Error(err.Error())
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		l.Error(err.Error())
	}
	MongoComponent.NewWriter(key, client)
	MongoComponent.NewReader(key, client)
	MongoNameMap[key] = dbname
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
