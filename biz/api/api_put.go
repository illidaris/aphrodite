package api

import (
	"net/http"

	"github.com/illidaris/aphrodite/dto"
	restSender "github.com/illidaris/rest/sender"
)

func NewPutAPI[Req IRequest, T any](param Req) *PutBaseAPI[Req, T] {
	return &PutBaseAPI[Req, T]{
		PostBaseAPI: PostBaseAPI[Req, T]{
			Request:  param,
			Response: new(dto.PtrResponse[T]),
		},
	}
}

var _ = restSender.IRequest(&PutBaseAPI[IRequest, any]{})

type PutBaseAPI[Req IRequest, T any] struct {
	PostBaseAPI[Req, T]
}

func (r PutBaseAPI[Req, T]) GetMethod() string {
	return http.MethodPut
}
