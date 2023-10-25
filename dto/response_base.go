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

func ErrorResponse(err error) *Response {
	res := &Response{}
	res.Code = -1
	res.Message = err.Error()
	return res
}

func OkResponse(data interface{}) *Response {
	res := &Response{}
	res.Message = "success"
	res.Data = data
	return res
}

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
