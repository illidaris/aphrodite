package dto

import "github.com/illidaris/aphrodite/pkg/exception"

type BaseResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	BaseResponse
	Data interface{} `json:"data"`
}

// ErrorResponse 函数接收一个错误对象 err，返回一个指向 Response 结构体的指针 res。
// 该函数用于生成一个错误响应对象，将错误信息赋值给 Response 结构体的 Message 字段，
// 并将 Code 字段设为 -1。
func ErrorResponse(err error) *Response {
	res := &Response{}
	res.Code = -1
	res.Message = err.Error()
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
			},
			Data: data,
		}
	}
	return OkResponse(data)
}
