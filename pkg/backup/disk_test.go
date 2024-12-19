package backup

import (
	"context"
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestB(t *testing.T) {
	mm := map[int]map[int]*Person{
		3: {
			1: {
				Name: "tes2222t",
				Age:  1,
			},
		},
		1: {
			1: {
				Name: "t2024",
				Age:  2024,
			},
		},
	}
	add(mm)
	bs, _ := json.Marshal(mm)
	println(string(bs))
}

type Person struct {
	Name string
	Age  int
}

func add(m map[int]map[int]*Person) {
	m[1] = map[int]*Person{}
	m[1] = map[int]*Person{
		1: {
			Name: "test",
			Age:  1,
		},
	}
}
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
