package convert

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestValuesRawEncode(t *testing.T) {
	convey.Convey("TestValuesRawEncode", t, func() {
		convey.Convey("empty", func() {
			v := url.Values{}
			v.Add("ticket", "@ml6sqYBGgTKmQNajnKNkaj8yksCAY++adIhlGIqfTiKyvBqOIkzdJ6WRgP+nO+wtVItqKbX4iZ+mFIYkyPJjpQ==")
			v.Add("timestamp", "1650941858")
			v.Add("nonce_str", "Wm3WZYTPz0wzccnW")
			bs := md5.Sum([]byte(ValuesRawEncode(v)))
			convey.So("3f7b739a91a52cb7d85c4f89c5f611fe", convey.ShouldEqual, fmt.Sprintf("%x", bs))
		})
	})
}
