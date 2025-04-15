package qrcodes

import (
	"encoding/base64"
	"testing"

	"github.com/skip2/go-qrcode"
	"github.com/smartystreets/goconvey/convey"
)

// TestWriteQrCode
func TestWriteQrCode(t *testing.T) {
	convey.Convey("WriteQrCode", t, func() {
		convey.Convey("WriteQrCode Success", func() {
			raw := "https://github.com/samber/do"
			bs, err := WriteQrCode(raw, qrcode.Medium, 256, "")
			convey.So(err, convey.ShouldBeNil)
			encodeStr := base64.StdEncoding.EncodeToString(bs)
			res, err := ParseQrCode(encodeStr)
			convey.So(err, convey.ShouldBeNil)
			convey.So(res, convey.ShouldEqual, raw)
		})
	})
}
