package backup

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestBackup(t *testing.T) {
	convey.Convey("TestBackup", t, func() {
		tempPath, err := os.MkdirTemp(os.TempDir(), "tmp")
		if err != nil {
			convey.So(err, convey.ShouldBeNil)
		}
		defer os.Remove(tempPath)

		convey.Convey("Save And Load", func() {
			ctx := context.Background()
			tempFile := path.Join(tempPath, "test.json")
			type Person struct {
				Name string
				Age  int
			}
			m := map[int]map[int]*Person{
				3: {
					1: {
						Name: "test",
						Age:  1,
					},
				},
				1: {
					222: {
						Name: "te2st",
						Age:  1,
					},
				},
			}
			err = DiskSave(ctx, tempFile, m)
			if err != nil {
				convey.So(err, convey.ShouldBeNil)
			}

			resp := map[int]map[int]*Person{}
			err := DiskLoad(ctx, tempFile, &resp)
			if err != nil {
				convey.So(err, convey.ShouldBeNil)
			}
			convey.So(resp, convey.ShouldEqual, m)
		})
	})
}
