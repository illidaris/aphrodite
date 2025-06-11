package idsnow

import (
	"time"

	"github.com/spf13/cast"
)

type Option func(*options)

func WithLenSequence(v int) Option {
	return func(opts *options) {
		opts.LenSequence = v
	}
}
func WithLenMachineID(v int) Option {
	return func(opts *options) {
		opts.LenMachineID = v
	}
}
func WithLenTimeUnix(v int) Option {
	return func(opts *options) {
		opts.LenTimeUnix = v
	}
}

func WithTimeUnix(v time.Duration) Option {
	return func(opts *options) {
		opts.TimeUnit = v
	}
}

func WithStartTime(v time.Time) Option {
	return func(opts *options) {
		opts.StartTime = v
	}
}

func WithNowFunc(f func() time.Time) Option {
	return func(opts *options) {
		opts.NowFunc = f
	}
}

func WithMachineID(f func() int) Option {
	return func(opts *options) {
		opts.MachineID = f
	}
}

func WithCheckMachineID(f func(int) bool) Option {
	return func(opts *options) {
		opts.CheckMachineID = f
	}
}
func newOptions(opts ...Option) options {
	options := defOptions() // 配置
	for _, opt := range opts {
		opt(&options)
	}
	diff := 63 - options.LenTotal()
	if diff > 0 {
		options.LenMachineID += diff
	}
	return options
}

func defOptions() options {
	opt := options{
		LenTimeUnix:  defaultBitsTime,
		LenSequence:  defaultBitsSequence,
		LenClock:     defaultBitsClock,
		LenMachineID: defaultBitsMachine,
		LenGene:      defaultBitGene,
		TimeUnit:     defaultTimeUnit,
		StartTime:    defaultStartTime,
		NowFunc:      func() time.Time { return time.Now() },
		MachineID:    nil, // TODO 默认方案：通过本机IP，在redis list中添加，同时返回机器IP对应的下标
		GeneFunc: func(key any, m int) int {
			if key == nil {
				return 0
			}
			return cast.ToInt(key) % m
		},
		CheckMachineID: nil,
	}
	return opt
}

type options struct {
	LenTimeUnix    int                // 长度-时间戳段
	LenSequence    int                // 长度-序列号段
	LenClock       int                // 长度-时钟段
	LenMachineID   int                // 长度-机器ID段
	LenGene        int                // 长度-基因段
	TimeUnit       time.Duration      // 时间戳-刻度单位 默认毫秒，不建议修改
	StartTime      time.Time          // 时间起点-默认2025年1月1日0点0分0秒，不建议修改
	NowFunc        func() time.Time   // 获取当前时间默认使用time.Now()
	MachineID      func() int         // 改造成实时计算(入参基因1（0~31），基因2（0~3）) 1-默认机器IP组成 2-自定义: NodeId（节点Id） 2^8 FrameId（主备帧数） 2^2 Gene2（基因2位取模） 2^2 Gene（基因4位取模） 2^4
	GeneFunc       func(any, int) int // 基因取模算法
	CheckMachineID func(int) bool     // 机器ID检查
}

func (o *options) LenTotal() int {
	return o.LenTimeUnix + o.LenClock + o.LenSequence + o.LenMachineID + o.LenGene
}
func (o *options) VaildOptions() error {
	// 检查配置的段总长
	if totalLen := o.LenTotal(); totalLen != 63 {
		return ErrTotalLength
	}
	// 限制序列长度
	if o.LenSequence < 0 || o.LenSequence > 16 {
		return ErrInvalidBitsSequence
	}
	// 限制机器码长度
	if o.LenMachineID < 0 || o.LenMachineID > 16 {
		return ErrInvalidBitsMachineID
	}
	// 确认时间戳序长度
	if o.LenTimeUnix < 32 {
		return ErrInvalidBitsTime
	}
	// 限制时间单位长度，不能小于1毫秒
	if o.TimeUnit < 0 || (o.TimeUnit > 0 && o.TimeUnit < time.Millisecond) {
		return ErrInvalidTimeUnit
	}
	// 限制起点时间
	if o.StartTime.After(time.Now()) {
		return ErrStartTimeAhead
	}
	// 确认机器Id是否合法
	if o.MachineID != nil {
		mid := o.MachineID()
		if mid < 0 || mid > 1<<o.LenMachineID-1 {
			return ErrInvalidMachineID
		}
	}
	return nil
}
func (o *options) LenSlice() []int {
	return []int{
		o.LenTimeUnix,  // 时间戳 长度 41
		o.LenClock,     // 时钟位 长度 1
		o.LenSequence,  // 序列位 长度 10
		o.LenMachineID, // 机器位 长度 7
		o.LenGene,      // 基因位 长度 4
	}
}

func (o *options) LenHeadSlice() []string {
	return []string{
		"时间戳(timestamp)", // 时间戳 长度 41
		"时钟位(clock)",     // 时钟位 长度 1
		"序列位(sequence)",  // 序列位 长度 10
		"机器位(machine)",   // 机器位 长度 7
		"基因位(gene)",      // 基因位 长度 4
	}
}

func (o *options) toId(vals ...int64) int64 {
	id := int64(0)
	for _, v := range IdPartsFrmVals(o.LenSlice(), vals...) {
		id |= v
	}
	return id
}

func (o *options) getUnit() int64 {
	return int64(o.TimeUnit)
}

func (o *options) getStartTimeUnix() int64 {
	return o.toInternalTime(o.StartTime)
}

func (o *options) toInternalTime(t time.Time) int64 {
	return t.UTC().UnixNano() / o.getUnit()
}

func (o *options) currentElapsedTime() int64 {
	return o.toInternalTime(o.NowFunc()) - o.getStartTimeUnix()
}

// sleep 休眠时间
func (o *options) sleep(overtime int64) {
	sleepTime := time.Duration(overtime*o.getUnit()) -
		time.Duration(time.Now().UTC().UnixNano()%o.getUnit())
	time.Sleep(sleepTime)
}
