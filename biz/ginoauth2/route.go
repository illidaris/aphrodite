package ginoauth2

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	apOAuth2 "github.com/illidaris/aphrodite/biz/oauth2"
)

func CookieMidleware(domain, secret string, seconds int) gin.HandlerFunc {
	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		Path:     "/",
		Domain:   domain,
		MaxAge:   seconds,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	return sessions.Sessions(apOAuth2.SESSION_KEY, store)
}
