package snowflake

import (
	"fmt"
	"strings"
	"time"

	"github.com/illidaris/aphrodite/pkg/convert"
)

func Compose(t time.Time, sequence, machineId, gene int, opts ...Option) (int64, error) {
	o := newOptions(opts...)
	elapsedTime := o.toInternalTime(t.UTC()) - o.getStartTimeUnix()
	if elapsedTime < 0 {
		return 0, ErrStartTimeAhead
	}
	if elapsedTime >= 1<<o.LenTimeUnix {
		return 0, ErrOverTimeLimit
	}
	if sequence < 0 || sequence >= 1<<o.LenSequence {
		return 0, ErrInvalidSequence
	}
	if machineId < 0 || machineId >= 1<<o.LenMachineID {
		return 0, ErrInvalidMachineID
	}
	if gene < 0 || gene >= 1<<o.LenGene {
		return 0, ErrInvalidGene
	}
	return o.toId(elapsedTime, int64(sequence), int64(machineId), int64(gene)), nil
}

// Decompose returns a set of Snowflake ID parts.
func Decompose(id int64, opts ...Option) []int64 {
	o := newOptions(opts...)
	return GetValsFrmId(o.LenSlice(), id)
}

func DecomposeStr(id int64, opts ...Option) string {
	o := newOptions(opts...) // 配置
	vals := GetValsFrmId(o.LenSlice(), id)
	headers := o.LenHeadSlice()
	strs := []string{fmt.Sprintf("产出的ID：%v", id)}
	for index, v := range vals {
		if index == 0 {
			strs = append(strs, fmt.Sprintf("%v:%v => %v", headers[index], v, convert.TimeFormat(o.ToTime(v))))
			continue
		}
		strs = append(strs, fmt.Sprintf("%v:%v", headers[index], v))
	}
	return strings.Join(strs, " ")
}
