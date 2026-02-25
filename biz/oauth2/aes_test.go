package oauth2

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAESEncode(t *testing.T) {
	secret := "0123456789abcdef"
	convey.Convey("TestAESEncode", t, func() {
		p := AuthorizeParam{
			Verifier: "12344444",
			Expire:   1234567890,
			BizId:    1,
		}
		encoded, err := AESEncode(p, secret)
		convey.So(err, convey.ShouldBeNil)

		decoded := AuthorizeParam{}
		err = AESDecode(&decoded, encoded, secret)
		convey.So(err, convey.ShouldBeNil)

		convey.So(decoded.Verifier, convey.ShouldEqual, p.Verifier)
		convey.So(decoded.Expire, convey.ShouldEqual, p.Expire)
		convey.So(decoded.BizId, convey.ShouldEqual, p.BizId)
	})
}
