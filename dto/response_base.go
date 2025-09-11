package dto

import "github.com/illidaris/aphrodite/pkg/exception"

// CheckResponseException 检查响应对象是否异常
// 参数：
// resp - 实现了IResponse接口的对象，表示一个响应实例
// 返回值：
// 返回一个exception.Exception类型的对象，表示检测到的异常或无异常（如果响应正常）
func CheckResponseException(resp IResponse) exception.Exception {
	// 检查resp是否为nil，如果是，则返回一个未响应的异常
	if resp == nil {
		return exception.ERR_UNRESPONSE.New("resp is nil")
	}
	// 如果resp不为nil，将其转换为异常对象并返回
	return resp.ToException()
}

type IResponse interface {
	ToException() exception.Exception
}

type BaseResponse struct {
	Code    int32  `json:"code"`
	SubCode int32  `json:"subCode"`
	Message string `json:"message"`
	Msg     string `json:"msg,omitempty"`
}

func (r BaseResponse) ToException() exception.Exception {
	if r.Code == 0 {
		return nil
	}
	exTp := exception.ExceptionType(r.Code)
	msg := r.Message
	if r.Msg != "" {
		msg = r.Msg
	}
	ex := exTp.New(msg)
	if r.SubCode != 0 {
		ex = ex.WithSubCode(r.SubCode)
	}
	return ex
}

type Response struct {
	BaseResponse
	Data interface{} `json:"data"`
}

type TResponse[T any] struct {
	BaseResponse
	Data T `json:"data"`
}

type PtrResponse[T any] struct {
	BaseResponse
	Data *T `json:"data"`
}

// ErrorResponse 函数接收一个错误对象 err，返回一个指向 Response 结构体的指针 res。
// 该函数用于生成一个错误响应对象，将错误信息赋值给 Response 结构体的 Message 字段，
// 并将 Code 字段设为 -1。
func ErrorResponse(err error) *Response {
	res := &Response{}
	res.Code = -1
	res.Message = err.Error()
	v, ok := err.(exception.Exception)
	if ok {
		res.Code = v.Code()
		res.SubCode = v.SubCode()
	}
	return res
}

// OkResponse函数用于生成一个成功的响应结果。
// 参数data为要返回的数据。
// 返回值为一个指向Response结构体的指针，其中Message字段为"success"，Data字段为data。
func OkResponse(data interface{}) *Response {
	res := &Response{}
	res.Message = "success"
	res.Data = data
	return res
}

// NewResponse 返回一个新的 Response 结构体指针。
// 如果 exception.Exception 不为 nil，则使用异常信息初始化 Response 结构体，
// 否则使用 data 初始化 OkResponse 结构体。
func NewResponse(data interface{}, ex exception.Exception) *Response {
	if ex != nil {
		return &Response{
			BaseResponse: BaseResponse{
				Code:    ex.Code(),
				Message: ex.Error(),
				SubCode: ex.SubCode(),
			},
			Data: data,
		}
	}
	return OkResponse(data)
}
