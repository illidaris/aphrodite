package api

import (
	"net/http"

	"github.com/illidaris/aphrodite/dto"
	restSender "github.com/illidaris/rest/sender"
)

func NewDeleteAPI[Req IRequest, T any](param Req) *DeleteBaseAPI[Req, T] {
	return &DeleteBaseAPI[Req, T]{
		GetBaseAPI: GetBaseAPI[Req, T]{
			Request:  param,
			Response: new(dto.TResponse[T]),
		},
	}
}

var _ = restSender.IRequest(&PutBaseAPI[IRequest, any]{})

type DeleteBaseAPI[Req IRequest, T any] struct {
	GetBaseAPI[Req, T]
}

func (r DeleteBaseAPI[Req, T]) GetMethod() string {
	return http.MethodDelete
}
