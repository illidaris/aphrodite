package idsnow

import "errors"

var (
	ErrTotalLength          = errors.New("total length must be 63")
	ErrInvalidBitsTime      = errors.New("bit length for time must be 32 or more")
	ErrInvalidBitsSequence  = errors.New("invalid bit length for sequence number")
	ErrInvalidBitsMachineID = errors.New("invalid bit length for machine id")
	ErrInvalidTimeUnit      = errors.New("invalid time unit")
	ErrInvalidSequence      = errors.New("invalid sequence number")
	ErrInvalidMachineID     = errors.New("invalid machine id")
	ErrStartTimeAhead       = errors.New("start time is ahead")
	ErrOverTimeLimit        = errors.New("over the time limit")
	ErrNoPrivateAddress     = errors.New("no private ip address")
)
