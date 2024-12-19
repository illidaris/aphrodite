package mongoex

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/illidaris/aphrodite/component/embedded"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	ErrCurNil      = errors.New("cur is nil")
	MongoComponent = embedded.NewComponent[*options.ClientOptions]()
	MongoNameMap   = map[string]string{}
	getKey         func(ctx context.Context) string
)

func SetGetKeyFunc(f func(ctx context.Context) string) {
	getKey = f
}

func CheckLink(ctx context.Context, name, uri string) error {
	clientOptions := options.Client().ApplyURI(uri).SetLoggerOptions(NewLoggerOptions())
	return Invoke(ctx, clientOptions, func(c *mongo.Client) error {
		subCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		return c.Ping(subCtx, readpref.Primary())
	})
}

func Invoke(ctx context.Context, opts *options.ClientOptions, cb func(c *mongo.Client) error) error {
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		println(err.Error())
	}
	defer func() {
		if disConErr := client.Disconnect(ctx); disConErr != nil {
			println(disConErr.Error())
		}
	}()
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("mongo_invoke_panic：%v \n", r)
		}
	}()
	return cb(client)
}

func NewMongo(key, dbname, conn string) {
	v := options.Client().ApplyURI(conn).SetLoggerOptions(NewLoggerOptions())
	MongoComponent.NewWriter(key, v)
	MongoComponent.NewReader(key, v)
	MongoNameMap[key] = dbname
}

// GetNamedMongoClient from mongo map
func GetNamedMongoClient(key string) *options.ClientOptions {
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
