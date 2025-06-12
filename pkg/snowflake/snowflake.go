package snowflake

import (
	"sync"
)

func NextIdFunc(opts ...Option) (func(key any) (int64, error), error) {
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
		return nil, err
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
			elapsedTime = current // 1. 进入下一个时间刻度，同时序列号从0开始
			sequence = 0
		} else {
			sequence = (sequence + 1) & maskSequence
			if sequence == 0 {
				elapsedTime++
				overtime := elapsedTime - current
				// TODO 等待时候超过一定时间，或者达到一定需求时，启用逻辑时钟, 后续代码中实现 逻辑时钟，原理参考：阿里Butterfly框架
				// TODO 逻辑时钟: 时间起点为进程开始时间点,序列号从0开始累计，到达限制之后，逻辑时累加一个单位时间，序列总0开始。
				options.sleep(overtime)
			}
		}
		// 时间超限
		if elapsedTime >= 1<<options.LenTimeUnix {
			return 0, ErrOverTimeLimit
		}
		return options.toId(
				elapsedTime,     // 相对时间戳
				clock,           // 时钟Id
				int64(sequence), // 序列Id
				int64(machine),  // 机器Id
				int64(gene)),    // 基因Id (由关键Id根据基因位长度取模生成)
			nil
	}, nil
}
