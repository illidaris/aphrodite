package idsnow

import (
	"context"
	"math/rand"

	"github.com/illidaris/aphrodite/idgenerate/dep"
	"github.com/illidaris/aphrodite/pkg/snowflake"
)

var _ = dep.IIDGenerate(&IdGenerate{})

func NewIdGenerate() *IdGenerate {
	return &IdGenerate{
		// Generaters:        generaters,
		// MchDstManager:     defaultMachineManager,
		// BackupMachineKeys: getMachineKeysOrInit,
	}
}

type IdGenerate struct {
	Generaters        []func(key any) (int64, error)
	MchDstManager     IMachineManager
	BackupMachineKeys func(ctx context.Context, dir string, num int, register func(string)) ([]string, error)
}

func (ig *IdGenerate) NewIDX(ctx context.Context, key string) int64 {
	id, _ := ig.NewID(ctx, key)
	return id
}
func (ig *IdGenerate) NewID(ctx context.Context, key string) (int64, error) {
	l := len(ig.Generaters)
	if l == 0 {
		return 0, snowflake.ErrHasNoGenerater
	}
	index := rand.Intn(l)
	id, err := ig.Generaters[index](key)
	return id, err
}
func (ig *IdGenerate) NewIDIterate(ctx context.Context, iterate func(int64), key string, opts ...dep.Option) error {
	panic("no impl")
}

func (ig *IdGenerate) Run(ctx context.Context, dir string, num int, opts ...snowflake.Option) error {
	if dir == "" {
		dir = "tmp"
	}
	if num < 1 || num > 128 {
		num = 4
	}
	mKeys, err := ig.BackupMachineKeys(ctx, dir, num, ig.MchDstManager.Register)
	if err != nil {
		return err
	}
	machineIdMap := ig.MchDstManager.GetMacgineIds()
	if ig.Generaters == nil {
		ig.Generaters = []func(key any) (int64, error){}
	}
	for _, mKey := range mKeys {
		mid, ok := machineIdMap[mKey]
		if ok {
			opts = append(opts, snowflake.WithMachineID(func() int {
				return mid
			}))
			f, err := snowflake.NextIdFunc(opts...)
			if err != nil {
				return err
			}
			ig.Generaters = append(ig.Generaters, f)
		}
	}
	return nil
}
