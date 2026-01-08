package cache

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/illidaris/aphrodite/pkg/dependency"
)

var _ = dependency.ILuaCache(redisTest{})

type redisTest struct {
	core *redis.Client
}

func (i redisTest) Eval(script string, keys []string, args ...any) (any, error) {
	return i.EvalContext(context.Background(), script, keys, args...)
}
func (i redisTest) EvalContext(ctx context.Context, script string, keys []string, args ...any) (any, error) {
	return i.core.Eval(ctx, script, keys, args...).Result()
}
func (i redisTest) Delete(key string) (int64, error) {
	return i.DeleteContext(context.Background(), key)
}
func (i redisTest) DeleteContext(ctx context.Context, key string) (int64, error) {
	return i.core.Del(ctx, key).Result()
}
func (i redisTest) Get(key string) (string, error) {
	return i.GetContext(context.Background(), key)
}
func (i redisTest) GetContext(ctx context.Context, key string) (string, error) {
	return i.core.Get(ctx, key).Result()
}

func TestNewsInfoForCache(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()

	// newsID := 123456789
	// key := fmt.Sprintf("news_redis_cache_%d", newsID)

	// mock ignoring `call api()`

	// mock.ExpectGet(key).RedisNil()
	// mock.Regexp().ExpectSet(key, `[a-z]+`, 30*time.Minute).SetErr(errors.New("FAIL"))
	mock.ExpectEval(LUA_ST_INC, []string{"_limiter:default:0:123"}).RedisNil()

	cache := redisTest{}
	cache.core = db

	res, err := LimitIncr(ctx, 0, "123", WithLimitCache(cache))
	if err != nil {
		t.Error("wrong error")
	}
	println(res)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
