package middleware

import (
	"context"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/illidaris/aphrodite/dto"
	"github.com/illidaris/aphrodite/pkg/check"
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
	WhiteIds     []string      // ID 白名单；为空时不参与组合判断
	WhiteIps     []string      // IP 白名单；为空时不参与组合判断
	WhiteHeaders []string      // 请求头白名单；为空时不参与组合判断
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

// IsInWhite 判断请求是否命中白名单。
// 仅参与判断已配置的白名单项：Or 策略下任一条件命中即放行，And 策略下需全部命中。
// 未配置任何白名单项或 CondType 非法时返回 false。
func (cfg SieveConfig) IsInWhite(id, ip, head string) bool {
	// 收集已配置的白名单条件，再按组合策略判断。
	checks := make([]WhiteCheck[string], 0, 3)
	// 空白名单不参与组合判断，避免未配置的条件影响 And/Or 结果。
	if len(cfg.WhiteIds) > 0 {
		checks = append(checks, WhiteCheck[string]{Value: id, Whites: cfg.WhiteIds})
	}
	if len(cfg.WhiteIps) > 0 {
		checks = append(checks, WhiteCheck[string]{Value: ip, Whites: cfg.WhiteIps, Check: func(s []string, e string) bool {
			return check.IsIPInWhitelist(e, s)
		}})
	}
	if len(cfg.WhiteHeaders) > 0 {
		checks = append(checks, WhiteCheck[string]{Value: head, Whites: cfg.WhiteHeaders})
	}
	// 没有任何已配置的白名单，视为未命中。
	if len(checks) == 0 {
		return false
	}
	switch cfg.CondType {
	case WhiteCondTypeAnd:
		for _, check := range checks {
			if !check.Valid() {
				return false
			}
		}
		return true
	case WhiteCondTypeOr:
		for _, check := range checks {
			if check.Valid() {
				return true
			}
		}
		return false
	default: // 未知策略按不命中处理，避免配置错误时意外放行请求。
		return false
	}
}

// WhiteCheck 描述单个白名单条件及其匹配规则。
// Whites 为空时该条件无效；Check 为空时使用 slices.Contains 判断是否命中。
type WhiteCheck[T comparable] struct {
	Value  T                 // 待校验的值
	Whites []T               // 白名单列表
	Check  func([]T, T) bool // 自定义匹配函数，参数依次为白名单列表和待校验值
}

// Valid 判断当前值是否命中白名单。
func (i WhiteCheck[T]) Valid() bool {
	// 未配置白名单时视为未命中。
	if len(i.Whites) == 0 {
		return false
	}
	// 未配置自定义匹配函数时，使用 slices.Contains 进行精确匹配。
	if i.Check != nil {
		return i.Check(i.Whites, i.Value)
	}
	return slices.Contains(i.Whites, i.Value)
}

// WhiteCondType 定义多个白名单条件之间的组合关系。
type WhiteCondType int8

const (
	// WhiteCondTypeOr 任意一个白名单条件命中即可放行。
	WhiteCondTypeOr WhiteCondType = iota
	// WhiteCondTypeAnd 所有白名单条件都命中才放行。
	WhiteCondTypeAnd
)
