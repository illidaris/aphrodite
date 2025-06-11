package idsnow

import (
	"context"
	"path"

	"github.com/google/uuid"
	"github.com/illidaris/aphrodite/pkg/backup"
	fileex "github.com/illidaris/file/path"
)

func GetMachineKeysOrInit(ctx context.Context, dir string, num int, register func(string)) ([]string, error) {
	machineKeys := []string{}
	if err := fileex.MkdirIfNotExist(dir); err != nil {
		return machineKeys, err
	}
	machineIdsFullFile := path.Join(dir, machineIdsFile)
	_ = backup.DiskLoad(ctx, machineIdsFullFile, &machineKeys)
	defer backup.DiskSave(ctx, machineIdsFullFile, &machineKeys)
	for i := 0; i < num-len(machineKeys); i++ {
		mKey := uuid.NewString()
		register(mKey)
		machineKeys = append(machineKeys, mKey)
	}
	return machineKeys, nil
}

func Offset(partLens []int, index int) int64 {
	l := len(partLens)
	if l == 0 {
		return 0
	}
	if index > l-1 {
		return 0
	}
	return offsetSum(partLens[index:]...)
}

func IdPartsFrmVals(partLens []int, vals ...int64) []int64 {
	parts := []int64{}
	for i := 0; i < len(partLens); i++ {
		val := int64(0)
		if i < len(vals) {
			val = vals[i]
		}
		offset := Offset(partLens, i+1)
		parts = append(parts, val<<offset)
	}
	return parts
}
func GetValsFrmId(partLens []int, id int64) []int64 {
	vals := []int64{}
	for i := 0; i < len(partLens); i++ {
		offset := Offset(partLens, i+1)
		maskSequence := (int64(1)<<int64(partLens[i]) - int64(1)) << offset
		vals = append(vals, (id&maskSequence)>>offset)
	}
	return vals
}

func offsetSum(ls ...int) int64 {
	t := int64(0)
	for _, l := range ls {
		t += int64(l)
	}
	return t
}
