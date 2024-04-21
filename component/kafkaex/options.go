package kafkaex

// OptionsFunc 是对 Options 结构体进行配置的函数类型。
type OptionsFunc func(o *Options)

// Options 结构体包含了连接信息所需的配置参数。
type Options struct {
	Addrs []string // 服务地址列表
	App   string   // 应用标识
	User  string   // 用户名
	Pwd   string   // 密码
}

// WithAddr 用于为 Options 添加一个或多个服务地址。
// 返回一个函数，该函数接受一个 Options 指针，将传入的地址添加到 Options 的 Addrs 列表中。
func WithAddr(addrs ...string) func(o *Options) {
	return func(o *Options) {
		if o.Addrs == nil {
			o.Addrs = []string{}
		}
		o.Addrs = append(o.Addrs, addrs...)
	}
}

// WithApp 用于设置 Options 的 App 字段。
// 返回一个函数，该函数接受一个 Options 指针，将传入的字符串赋值给 Options 的 App 字段。
func WithApp(app string) func(o *Options) {
	return func(o *Options) {
		o.App = app
	}
}

// WithUser 用于设置 Options 的 User 字段。
// 返回一个函数，该函数接受一个 Options 指针，将传入的字符串赋值给 Options 的 User 字段。
func WithUser(user string) func(o *Options) {
	return func(o *Options) {
		o.User = user
	}
}

// WithPwd 用于设置 Options 的 Pwd 字段。
// 返回一个函数，该函数接受一个 Options 指针，将传入的字符串赋值给 Options 的 Pwd 字段。
func WithPwd(pwd string) func(o *Options) {
	return func(o *Options) {
		o.Pwd = pwd
	}
}
