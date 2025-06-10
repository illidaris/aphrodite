package idsnow

import (
	"sync"
)

func NextIdFunc(opts ...Option) func(key any) (int64, error) {
	options := newOptions(opts...) // 配置
	var (
		mutex       = new(sync.Mutex)            // 锁
		elapsedTime int64                        // 上一个Id的时间戳
		machine     int                          // 机器ID
		sequence    = 1<<options.LenSequence - 1 // 当前序列ID
		clock       int64
	)
	err := options.VaildOptions()
	if err != nil {
		return func(key any) (int64, error) { return 0, err }
	}
	if options.MachineID != nil {
		machine = options.MachineID()
	}
	return func(key any) (int64, error) {
		maskSequence := 1<<options.LenSequence - 1 // 构建【序列段】
		gene := options.GeneFunc(key, 1<<options.LenGene)
		mutex.Lock()                            // 加锁
		defer mutex.Unlock()                    // 解锁
		current := options.currentElapsedTime() // 当前偏移时间戳
		if elapsedTime < current {              // 当前偏移时间戳 大于 历史偏移时间戳
			// 1. 进入下一个时间刻度，同时序列号从0开始
			elapsedTime = current
			sequence = 0
		} else {
			// TODO: 处理时间回拨，添加历史时钟
			sequence = (sequence + 1) & maskSequence
			if sequence == 0 {
				elapsedTime++
				overtime := elapsedTime - current
				options.sleep(overtime)
			}
		}
		// 时间超限
		if elapsedTime >= 1<<options.LenTimeUnix {
			return 0, ErrOverTimeLimit
		}
		return options.toId(
				elapsedTime,
				clock,
				int64(sequence),
				int64(machine),
				int64(gene)),
			nil
	}
}
