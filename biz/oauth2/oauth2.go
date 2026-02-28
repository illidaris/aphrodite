package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/illidaris/aphrodite/pkg/contextex"
	"github.com/oklog/ulid/v2"
	"golang.org/x/oauth2"
)

// GetAuthorizeURl 获取授权URL
func GetAuthorizeURl(ctx context.Context, opts ...Option) (string, AuthorizeParam, string, error) {
	opt := NewOptions(opts...)
	state := ulid.Make().String()
	verifier := oauth2.GenerateVerifier()
	dur := opt.GetCodeChallengeExpireHandle(ctx)
	param := AuthorizeParam{
		Verifier: verifier,
		BizId:    contextex.GetBizId(ctx),
		Expire:   time.Now().Add(dur).Unix(),
	}
	if opt.Cache != nil {
		key := fmt.Sprintf(CACHE_KEY_CODE_VERIFIER, state)
		valBs, _ := json.Marshal(param)
		err := opt.Cache.SetCtx(ctx, key, string(valBs), dur)
		if err != nil {
			return "", param, param.Encode(opt.GetBusiSecretHandle(ctx)), err
		}
	}
	url := opt.GetOAuth2Config(ctx).AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))
	return url, param, param.Encode(opt.GetBusiSecretHandle(ctx)), nil
}

// OAuthCallback 处理OAuth回调
func OAuthCallback(ctx context.Context, param *OAuthCallbackParam, findCodeVerifier func(string) string, opts ...Option) (*oauth2.Token, error) {
	opt := NewOptions(opts...)
	bizId := contextex.GetBizId(ctx)
	cacheParam := &AuthorizeParam{}
	if findCodeVerifier != nil {
		v := findCodeVerifier(param.State)
		err := cacheParam.Decode(v, opt.GetBusiSecretHandle(ctx))
		if err != nil {
			return nil, err
		}
	} else if opt.Cache != nil {
		key := fmt.Sprintf(CACHE_KEY_CODE_VERIFIER, param.State)
		res, err := opt.Cache.GetCtx(ctx, key)
		if err != nil {
			return nil, err
		}
		defer opt.Cache.DelCtx(ctx, key)
		err = json.Unmarshal([]byte(res), cacheParam)
		if err != nil {
			return nil, err
		}
	}
	if err := cacheParam.Valid(bizId); err != nil {
		return nil, err
	}
	return opt.GetOAuth2Config(ctx).Exchange(ctx, param.Code, oauth2.VerifierOption(cacheParam.Verifier))
}
