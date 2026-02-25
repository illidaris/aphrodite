package oauth2

import (
	"context"
	"sync"
	"time"
)

const (
	CACHE_KEY_CODE_VERIFIER = "_aph_oauth2:code_verifier:%d:%s"
	SESSION_KEY             = "_aph_oauth2_session"
)

type ICache interface {
	SetCtx(context.Context, string, string, time.Duration) error
	GetCtx(context.Context, string) (string, error)
}

type defaultCache struct {
	*sync.Map
}

func (c defaultCache) SetCtx(ctx context.Context, key, value string, duration time.Duration) error {
	c.Store(key, value)
	return nil
}
func (c defaultCache) GetCtx(ctx context.Context, key string) (string, error) {
	if value, ok := c.Load(key); ok {
		return value.(string), nil
	}
	return "", nil
}
