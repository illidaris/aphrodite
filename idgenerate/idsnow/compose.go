package idsnow

import (
	"strings"
	"time"

	"github.com/spf13/cast"
)

// ToTime returns the time when the given ID was generated.
func (o *options) ToTime(id int64) time.Time {
	return time.Unix(0, (o.getStartTimeUnix()+o.timePart(id))*o.getUnit())
}

// Compose creates a Snowflake ID from its components.
// The time parameter should be the time when the ID was generated.
// The sequence parameter should be between 0 and 2^BitsSequence-1 (inclusive).
// The machineID parameter should be between 0 and 2^BitsMachineID-1 (inclusive).
func (o *options) Compose(t time.Time, sequence, machineID int) (int64, error) {
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

	if machineID < 0 || machineID >= 1<<o.LenMachineID {
		return 0, ErrInvalidMachineID
	}

	return elapsedTime<<(int64(o.LenSequence)+int64(o.LenMachineID)) |
		int64(sequence)<<int64(o.LenMachineID) |
		int64(machineID), nil
}

// Decompose returns a set of Snowflake ID parts.
func (o *options) Decompose(id int64) map[string]int64 {
	time := o.timePart(id)
	sequence := o.sequencePart(id)
	machine := o.machinePart(id)
	return map[string]int64{
		"id":       id,
		"time":     time,
		"sequence": sequence,
		"machine":  machine,
	}
}

func DecomposeStr(id int64, opts ...Option) string {
	o := newOptions(opts...) // 配置
	vals := o.GetValsFrmId(id)
	strs := []string{cast.ToString(id)}
	for _, v := range vals {
		strs = append(strs, cast.ToString(v))
	}
	return strings.Join(strs, " ")
}

func (o *options) timePart(id int64) int64 {
	return id >> (int64(o.LenSequence) + int64(o.LenMachineID))
}

func (o *options) sequencePart(id int64) int64 {
	maskSequence := int64((1<<o.LenSequence - 1) << o.LenMachineID)
	return (id & maskSequence) >> o.LenMachineID
}

func (o *options) machinePart(id int64) int64 {
	maskMachine := int64(1<<o.LenMachineID - 1)
	return id & maskMachine
}
