package kafkaex

import "errors"

const (
	ERR_CONSUMER_HASRUN  = "consumer has run"
	ERR_CONSUMER_NOTRUN  = "consumer no run"
	ERR_CONSUMER_EXIST   = "consumer exist"
	ERR_CONSUMER_NOFOUND = "consumer no found"
	ERR_GROUP_NOFOUND    = "group no found"
)

var (
	ErrConsumerHasRun  = errors.New(ERR_CONSUMER_HASRUN)
	ErrConsumerNotRun  = errors.New(ERR_CONSUMER_NOTRUN)
	ErrConsumerExist   = errors.New(ERR_CONSUMER_EXIST)
	ErrConsumerNoFound = errors.New(ERR_CONSUMER_NOFOUND)
	ErrGroupNoFound    = errors.New(ERR_GROUP_NOFOUND)
)
