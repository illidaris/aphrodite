package idsnow

import (
	"time"
)

const (
	defaultTimeUnit     = 1e6 // 10^6 毫秒
	defaultBitsTime     = 41  // 默认时间长度
	defaultBitsSequence = 10  // 默认序列长度
	defaultBitsClock    = 1   // 时钟长度 0-默认机器时钟 1-自定义时钟（逻辑时钟）
	defaultBitsMachine  = 7   // 默认机器ID长度
	defaultBitGene      = 4   // 默认基因长度
)

var (
	defaultStartTime = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
)
