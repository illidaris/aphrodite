package ginoauth2

import (
	"context"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	apOAuth2 "github.com/illidaris/aphrodite/biz/oauth2"
	"github.com/illidaris/aphrodite/dto"
	"github.com/illidaris/aphrodite/pkg/exception"
	"github.com/spf13/cast"
	"golang.org/x/oauth2"
)

const (
	LOGIN_BY_OAUTH2 = "_aph_by_oauth2"
)

func LoginByOAuth2RedirectController(opts ...apOAuth2.Option) func(c *gin.Context) {
	return func(c *gin.Context) {
		// 获取默认Session
		session := sessions.Default(c)
		// 设置Session中的键值对
		url, _, str, err := apOAuth2.GetAuthorizeURl(c.Request.Context(), opts...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewResponse(nil, exception.ERR_COMMON_BADPARAM.Wrap(err)))
			return
		}
		session.Set(LOGIN_BY_OAUTH2, str)
		session.Save()
		// 重定向到授权URL
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func LoginByOAuth2Controller(opts ...apOAuth2.Option) func(c *gin.Context) {
	return func(c *gin.Context) {
		// 获取默认Session
		session := sessions.Default(c)
		// 设置Session中的键值对
		url, _, str, err := apOAuth2.GetAuthorizeURl(c.Request.Context(), opts...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewResponse(nil, exception.ERR_COMMON_BADPARAM.Wrap(err)))
			return
		}
		session.Set(LOGIN_BY_OAUTH2, str)
		session.Save()
		// 重定向到授权URL
		c.JSON(http.StatusOK, dto.NewResponse(url, nil))
	}
}

func CallbackOAuth2Controller(handle func(ctx context.Context, token *oauth2.Token) (any, error), opts ...apOAuth2.Option) func(c *gin.Context) {
	return func(c *gin.Context) {
		param := &apOAuth2.OAuthCallbackParam{}
		if err := c.ShouldBind(param); err != nil {
			c.JSON(http.StatusOK, dto.NewResponse(nil, exception.ERR_COMMON_BADPARAM.Wrap(err)))
			return
		}
		session := sessions.Default(c)
		value := session.Get(LOGIN_BY_OAUTH2)
		defer session.Delete(LOGIN_BY_OAUTH2)
		token, err := apOAuth2.OAuthCallback(c.Request.Context(), param, func(state string) string {
			return cast.ToString(value)
		}, opts...)
		if err != nil {
			c.JSON(http.StatusOK, dto.NewResponse(nil, exception.ERR_COMMON_UNAUTH.Wrap(err)))
			return
		}
		if handle == nil {
			c.JSON(http.StatusOK, dto.NewResponse(nil, exception.ERR_COMMON_UNAUTH.New("没有验证函数")))
			return
		}
		resp, err := handle(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusOK, dto.NewResponse(resp, exception.ERR_COMMON_UNAUTH.Wrap(err)))
			return
		}
		c.JSON(http.StatusOK, dto.NewResponse(resp, nil))
	}
}
