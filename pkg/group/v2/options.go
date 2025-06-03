package v2

const (
	// DefaultBatchSize 默认的批次大小
	DEF_BATCH_SIZE      = 100
	DEF_MAX_GROUP_COUNT = 10000
)

type options struct {
	ParallelismMax int
	Batch          int
}

func (o options) GetRealBatch(srcsLen int) int {
	if o.Batch == 0 {
		return 1
	}
	// 没有数据
	if srcsLen == 0 {
		return 1
	}
	// 防止负提升
	if srcsLen < o.Batch {
		return 1
	}
	return o.Batch
}

func (o options) GetRealCount(srcsLen, count int) (int, int) {
	realBatch := o.GetRealBatch(srcsLen)
	if count > o.ParallelismMax {
		if srcsLen%o.ParallelismMax == 0 {
			realBatch = srcsLen / o.ParallelismMax
		} else {
			realBatch = srcsLen/o.ParallelismMax + 1
		}
		return realBatch, o.ParallelismMax
	}
	return realBatch, count
}

func newOptions(opts ...Option) *options {
	opt := &options{
		ParallelismMax: DEF_MAX_GROUP_COUNT,
		Batch:          DEF_BATCH_SIZE,
	}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

type Option func(*options)

func WithParallelismMax(parallelism int) Option {
	return func(o *options) {
		o.ParallelismMax = parallelism
	}
}

func WithBatch(batch int) Option {
	return func(o *options) {
		o.Batch = batch
	}
}
