package dep

type Option func(*options)

func NewOptions(opts ...Option) *options {
	opt := &options{
		Num: 1,
	}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

type options struct {
	Num int64
}

func WithNum(num int64) Option {
	return func(o *options) {
		o.Num = num
	}
}
