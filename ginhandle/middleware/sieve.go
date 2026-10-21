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
		// 未启用时不进行拦截。
		if !cfg.Enable {
			c.Next()
			return
		}
		id := SieveWhiteIdFrmCtx(ctx)
		ip := c.ClientIP()
		headVal := c.Request.Header.Get(HEAD_WHITE)
		if cfg.IsInWhite(id, ip, headVal) {
			c.Next()
			return
		}
		// 时间配置使用 Unix 秒；尚未开始或已经过期时返回 403。
		now := when.Now(ctx)
		if cfg.NotBefore > 0 && now.Unix() < int64(cfg.NotBefore) {
			c.AbortWithStatusJSON(403, dto.NewResponse(nil, exception.ERR_BUSI.New("Forbidden Nbf")))
			return
		}
		if cfg.Expire > 0 && now.Unix() > int64(cfg.Expire) {
			c.AbortWithStatusJSON(403, dto.NewResponse(nil, exception.ERR_BUSI.New("Forbidden Exp")))
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
	WhiteIds     []string      // ID 白名单；为空时不参与组合判断
	WhiteIps     []string      // IP 白名单；为空时不参与组合判断
	WhiteHeaders []string      // 请求头白名单；为空时不参与组合判断
}

// IsInWhite 按配置的组合策略判断请求是否命中白名单；未配置的白名单不参与判断。
func (cfg SieveConfig) IsInWhite(id, ip, head string) bool {
	configured := 0
	matched := false
	check := func(values []string, value string) bool {
		if len(values) == 0 {
			return true
		}
		configured++
		ok := slices.Contains(values, value)
		matched = matched || ok
		return ok
	}

	switch cfg.CondType {
	case WhiteCondTypeAnd:
		return check(cfg.WhiteIds, id) && check(cfg.WhiteIps, ip) && check(cfg.WhiteHeaders, head) && configured > 0
	case WhiteCondTypeOr:
		check(cfg.WhiteIds, id)
		check(cfg.WhiteIps, ip)
		check(cfg.WhiteHeaders, head)
		return configured > 0 && matched
	default: // 未知策略按不命中处理，避免配置错误意外放行请求。
		return false
	}
}

// IsInWhiteId 判断指定 ID 是否命中白名单；未配置白名单时视为匹配。
func (cfg SieveConfig) IsInWhiteId(id string) bool {
	if len(cfg.WhiteIds) == 0 {
		return true
	}
	return slices.Contains(cfg.WhiteIds, id)
}

// IsInWhiteIp 判断指定 IP 是否命中白名单；未配置白名单时视为匹配。
func (cfg SieveConfig) IsInWhiteIp(ip string) bool {
	if len(cfg.WhiteIps) == 0 {
		return true
	}
	return slices.Contains(cfg.WhiteIps, ip)
}

// IsInWhiteHeader 判断指定请求头值是否命中白名单；未配置白名单时视为匹配。
func (cfg SieveConfig) IsInWhiteHeader(header string) bool {
	if len(cfg.WhiteHeaders) == 0 {
		return true
	}
	return slices.Contains(cfg.WhiteHeaders, header)
}

// WhiteCondType 定义多个白名单条件之间的组合关系。
type WhiteCondType int8

const (
	// WhiteCondTypeOr 任意一个白名单条件命中即可放行。
	WhiteCondTypeOr WhiteCondType = iota
	// WhiteCondTypeAnd 所有白名单条件都命中才放行。
	WhiteCondTypeAnd
)
