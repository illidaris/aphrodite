package oauth2

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestGetAuthorizeURl(t *testing.T) {
	feishuAuth := "https://accounts.feishu.cn/open-apis/authen/v1/authorize"
	feishuToken := "https://open.feishu.cn/open-apis/authen/v2/oauth/token"
	callback := "https://www.illidaris.com/callback"
	ctx := context.Background()

	convey.Convey("TestGetAuthorizeURl", t, func() {
		url, _, _, err := GetAuthorizeURl(ctx,
			WithGetAuthUrlHandle(func(context.Context) string {
				return feishuAuth
			}),
			WithGetTokenUrlHandle(func(ctx context.Context) string {
				return feishuToken
			}),
			WithGetRedirectUrlHandle(func(ctx context.Context) string {
				return callback
			}),
			WithGetClientIdHandle(func(ctx context.Context) string {
				return "test_client_id"
			}),
			WithGetClientSecretHandle(func(ctx context.Context) string {
				return "test_client_secret"
			}),
			WithGetBusiSecretHandle(func(ctx context.Context) string {
				return "test_busi_secret"
			}),
		)
		println(url)
		//target := `https://accounts.feishu.cn/open-apis/authen/v1/authorize?client_id=test_client_id&code_challenge=mBfU8bWCfw6dKw3TwAI4mI-_ikwuDnSXBYBdvKmIbZk&code_challenge_method=S256&redirect_uri=https%3A%2F%2Fwww.illidaris.com%2Fcallback&response_type=code&state=fKdJ0K6IO5LkXYpH3ZgpL9vcOrCt9zyQjuwul%2F3cOMp5zVZaG5kd3Nbxl%2BrDXYsUADro4XcQS6FRIfFgu%2FKwmrTirrukDg%3D%3D`
		convey.So(err, convey.ShouldBeNil)
		//convey.So(url, convey.ShouldEqual, target)
	})
}
