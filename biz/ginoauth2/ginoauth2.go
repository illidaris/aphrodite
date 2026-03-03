package ginoauth2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	apOAuth2 "github.com/illidaris/aphrodite/biz/oauth2"
	"github.com/illidaris/aphrodite/dto"
	"github.com/illidaris/aphrodite/pkg/exception"
	"golang.org/x/oauth2"
)

func LoginByOAuth2RedirectController(opts ...apOAuth2.Option) func(c *gin.Context) {
	return func(c *gin.Context) {
		url, _, _, err := apOAuth2.GetAuthorizeURl(c.Request.Context(), opts...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewResponse(nil, exception.ERR_COMMON_BADPARAM.Wrap(err)))
			return
		}
		// 重定向到授权URL
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func LoginByOAuth2Controller(opts ...apOAuth2.Option) func(c *gin.Context) {
	return func(c *gin.Context) {
		url, _, _, err := apOAuth2.GetAuthorizeURl(c.Request.Context(), opts...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewResponse(nil, exception.ERR_COMMON_BADPARAM.Wrap(err)))
			return
		}
		// 重定向到授权URL
		c.JSON(http.StatusOK, dto.NewResponse(url, nil))
	}
}

func CallbackOAuth2Controller(handle func(ctx context.Context, token *oauth2.Token) (any, exception.Exception), opts ...apOAuth2.Option) func(c *gin.Context) {
	return func(c *gin.Context) {
		param := &apOAuth2.OAuthCallbackParam{}
		if err := c.ShouldBind(param); err != nil {
			c.JSON(http.StatusUnauthorized, dto.NewResponse(nil, exception.ERR_COMMON_BADPARAM.Wrap(fmt.Errorf("请求飞书登录态参数错误%v", err))))
			return
		}
		token, err := apOAuth2.OAuthCallback(c.Request.Context(), param, nil, opts...)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.NewResponse(nil, exception.ERR_COMMON_UNAUTH.Wrap(fmt.Errorf("获取飞书登录态失败%v", err))))
			return
		}
		if handle == nil {
			c.JSON(http.StatusUnauthorized, dto.NewResponse(nil, exception.ERR_COMMON_UNAUTH.New("没有配置验证过程")))
			return
		}
		resp, ex := handle(c.Request.Context(), token)
		if ex != nil {
			c.JSON(http.StatusOK, dto.NewResponse(resp, ex))
			return
		}
		c.JSON(http.StatusOK, dto.NewResponse(resp, nil))
	}
}
