package middleware

import (
	"context"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/illidaris/aphrodite/dto"
	"github.com/illidaris/aphrodite/pkg/exception"
	"github.com/illidaris/aphrodite/pkg/when"
)

// SieveMiddleware 根据业务配置、白名单和启用时间范围控制请求是否放行。
func SieveMiddleware(busi string, newConfig NewConfigHandle) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		cfg := newConfig(ctx, busi)
		ip := c.ClientIP()
		headVal := c.Request.Header.Get(HEAD_WHITE)
		if ex := cfg.Check(ctx, ip, headVal); ex != nil {
			c.AbortWithStatusJSON(403, dto.NewResponse(nil, ex))
			return
		}
		c.Next()
	}
}

// HEAD_WHITE 是用于请求头白名单校验的请求头名称。
const (
	HEAD_WHITE = "X-Aphrodite-SieveWhite"
)

type busiWhiteIdCtxKey struct{}

// WithSieveWhiteId 将业务 ID 白名单值写入上下文。
func WithSieveWhiteId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, busiWhiteIdCtxKey{}, id)
}

// SieveWhiteIdFrmCtx 从上下文读取业务 ID 白名单值；未设置时返回空字符串。
func SieveWhiteIdFrmCtx(ctx context.Context) string {
	id := ctx.Value(busiWhiteIdCtxKey{})
	if id == nil {
		return ""
	}
	value, ok := id.(string)
	if !ok {
		return ""
	}
	return value
}

// NewConfigHandle 根据请求上下文和业务标识生成筛选配置。
type NewConfigHandle func(context.Context, string) SieveConfig

// SieveConfig 描述筛选开关、有效时间范围及各类白名单。
type SieveConfig struct {
	Enable       bool          // 是否启用筛选
	NotBefore    int           // 开始生效时间（Unix 秒，0 表示不限制）
	Expire       int           // 失效时间（Unix 秒，0 表示不限制）
	CondType     WhiteCondType // 白名单组合策略：或/与
	WhiteIds     []string      // ID 白名单；为空时视为匹配
	WhiteIps     []string      // IP 白名单；为空时视为匹配
	WhiteHeaders []string      // 请求头白名单；为空时视为匹配
}

// Check 检查筛选配置是否允许当前请求通过。
func (cfg SieveConfig) Check(ctx context.Context, ip, head string) exception.Exception {
	// 未启用时不进行拦截。
	if !cfg.Enable {
		return nil
	}
	if cfg.IsInWhite(SieveWhiteIdFrmCtx(ctx), ip, head) {
		return nil
	}
	// 时间配置使用 Unix 秒；尚未开始或已经过期时返回 403 错误。
	now := when.Now(ctx)
	if cfg.NotBefore > 0 && now.Unix() < int64(cfg.NotBefore) {
		return exception.ERR_BUSI.New("Forbidden Nbf")
	}
	if cfg.Expire > 0 && now.Unix() > int64(cfg.Expire) {
		return exception.ERR_BUSI.New("Forbidden Exp")
	}
	return nil
}

// IsInWhite 按配置的组合策略组合各项白名单判断结果。
func (cfg SieveConfig) IsInWhite(id, ip, head string) bool {
	// 未配置任何白名单项时不额外限制；否则按策略合并已配置的条件。
	type whiteCheck struct{ configured, matched bool }
	checks := make([]whiteCheck, 0, 3)
	isCheck, matched := cfg.IsInWhiteId(id)
	checks = append(checks, whiteCheck{isCheck, matched})
	isCheck, matched = cfg.IsInWhiteIp(ip)
	checks = append(checks, whiteCheck{isCheck, matched})
	isCheck, matched = cfg.IsInWhiteHeader(head)
	checks = append(checks, whiteCheck{isCheck, matched})

	configured := false
	result := cfg.CondType == WhiteCondTypeAnd
	if cfg.CondType != WhiteCondTypeAnd && cfg.CondType != WhiteCondTypeOr {
		return false
	}
	for _, check := range checks {
		if !check.configured {
			continue
		}
		configured = true
		if cfg.CondType == WhiteCondTypeAnd {
			result = result && check.matched
		} else {
			result = result || check.matched
		}
	}
	if !configured {
		return true
	}
	return result
}

// IsInWhiteId 判断指定 ID 是否命中白名单；未配置白名单时视为匹配。
func (cfg SieveConfig) IsInWhiteId(id string) (bool, bool) {
	return hasVal(id, cfg.WhiteIds...)
}

// IsInWhiteIp 判断指定 IP 是否命中白名单；未配置白名单时视为匹配。
func (cfg SieveConfig) IsInWhiteIp(ip string) (bool, bool) {
	return hasVal(ip, cfg.WhiteIps...)
}

// IsInWhiteHeader 判断指定请求头值是否命中白名单；未配置白名单时视为匹配。
func (cfg SieveConfig) IsInWhiteHeader(header string) (bool, bool) {
	return hasVal(header, cfg.WhiteHeaders...)
}

func hasVal(v string, args ...string) (bool, bool) {
	isCheck, has := false, false
	l := len(args)
	isCheck = l > 0
	if l == 0 {
		return isCheck, has
	}
	if v == "" {
		return isCheck, false
	}
	return isCheck, slices.Contains(args, v)
}

// WhiteCondType 定义多个白名单条件之间的组合关系。
type WhiteCondType int8

const (
	// WhiteCondTypeOr 任意一个白名单条件命中即可放行。
	WhiteCondTypeOr WhiteCondType = iota
	// WhiteCondTypeAnd 所有白名单条件都命中才放行。
	WhiteCondTypeAnd
)
