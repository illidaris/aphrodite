package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// 定义一个函数类型，用于设置代理选项。
type ProxyOption func(*ProxyOptions)

// 定义代理选项的结构体，包含各种钩子函数。
type ProxyOptions struct {
	ErrHook      func(http.ResponseWriter, *http.Request, error) // 错误处理钩子，当代理过程中出现错误时调用。
	RequestHooks []func(*http.Request)                           // 请求处理钩子，用于在发送请求之前对请求进行修改或处理。
	ResponseHook func(*http.Response) error                      // 响应处理钩子，用于在返回响应之前对响应进行修改或处理。
	RewriteHook  func(pr *httputil.ProxyRequest)                 // 重写处理钩子，用于修改请求的URL或其他请求头信息。
}

// 使用错误处理钩子函数设置代理选项。
func WithErrHook(f func(http.ResponseWriter, *http.Request, error)) ProxyOption {
	return func(opts *ProxyOptions) {
		opts.ErrHook = f
	}
}

// 使用请求处理钩子函数设置代理选项。
func WithRequestHooks(f ...func(*http.Request)) ProxyOption {
	return func(opts *ProxyOptions) {
		opts.RequestHooks = f
	}
}

// 使用重写处理钩子函数设置代理选项。
func WithRewriteHook(f func(pr *httputil.ProxyRequest)) ProxyOption {
	return func(opts *ProxyOptions) {
		opts.RewriteHook = f
	}
}

// 使用响应处理钩子函数设置代理选项。
func WithResponseHook(f func(*http.Response) error) ProxyOption {
	return func(opts *ProxyOptions) {
		opts.ResponseHook = f
	}
}

// NewProxy 创建一个新的反向代理，支持通过选项进行定制。
// targetHost: 目标主机的URL。
// opt: 一个或多个代理选项函数，用于定制代理的行为。
// 返回值: 创建的反向代理实例和可能的错误。
func NewProxy(targetHost string, opt ...ProxyOption) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	option := ProxyOptions{}
	for _, o := range opt {
		o(&option)
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		for _, hook := range option.RequestHooks {
			hook(req)
		}
	}
	if option.RewriteHook != nil {
		proxy.Rewrite = option.RewriteHook
	}
	if option.ResponseHook != nil {
		proxy.ModifyResponse = option.ResponseHook
	}
	if option.ErrHook != nil {
		proxy.ErrorHandler = option.ErrHook
	}
	return proxy, nil
}

// ProxyRequestHandler 创建一个处理HTTP请求的函数，该函数使用指定的反向代理来处理请求。
// proxy: 用于处理请求的反向代理实例。
// 返回值: 一个处理HTTP请求的函数。
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}
