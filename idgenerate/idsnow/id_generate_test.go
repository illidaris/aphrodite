package idsnow

import (
	"context"
	"math/rand"
	"os"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cast"
)

func TestIdGenerate(t *testing.T) {
	convey.Convey("TestIdGenerate", t, func() {
		ctx := context.Background()
		mm := &testMachineManager{}
		mm.Register("test1")
		mm.Register("test2")
		mm.Register("test3")
		mm.Register("test4")

		idger := &IdGenerate{
			MchDstManager:     mm,
			BackupMachineKeys: GetMachineKeysOrInit,
		}

		dir, err := os.MkdirTemp("", "idsnow_test")
		convey.So(err, convey.ShouldBeNil)
		defer os.RemoveAll(dir)

		err = idger.Run(ctx, dir, 3)
		convey.So(err, convey.ShouldBeNil)

		convey.Convey("success", func() {
			for i := 0; i < 100; i++ {
				gene := rand.Int63n(123456789)
				md := gene % 16
				id, err := idger.NewID(ctx, cast.ToString(gene))
				convey.So(err, convey.ShouldBeNil)
				vals := Decompose(id)
				convey.So(vals[2], convey.ShouldBeBetweenOrEqual, 0, 1023)
				convey.So(vals[3], convey.ShouldBeBetween, 3, 8)
				convey.So(vals[4], convey.ShouldEqual, md)
			}
		})
	})
}

var _ = IMachineManager(&testMachineManager{})

type testMachineManager struct {
	s []string
	m map[string]int
}

func (t *testMachineManager) Register(id string) {
	if t.s == nil {
		t.s = make([]string, 0)
		t.m = make(map[string]int)
	}
	_, ok := t.m[id]
	if ok {
		return
	}
	t.s = append(t.s, id)
	t.m[id] = len(t.s) - 1
}
func (t *testMachineManager) GetMacgineIds() map[string]int {
	return t.m
}
