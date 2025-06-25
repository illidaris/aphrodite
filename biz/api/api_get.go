package api

import (
	"encoding/json"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/illidaris/aphrodite/dto"
	restSender "github.com/illidaris/rest/sender"
)

func NewGetAPI[Req IRequest, T any](param Req) *GetBaseAPI[Req, T] {
	return &GetBaseAPI[Req, T]{
		Request:  param,
		Response: new(dto.PtrResponse[T]),
	}
}

var _ = restSender.IRequest(&GetBaseAPI[IRequest, any]{})

type GetBaseAPI[Req IRequest, T any] struct {
	restSender.GETRequest `json:"-"`
	Request               Req
	Response              *dto.PtrResponse[T] `json:"-"`
}

func (r GetBaseAPI[Req, T]) GetUrlQuery() url.Values {
	u, _ := query.Values(r.Request)
	return u
}

func (r GetBaseAPI[Req, T]) Encode() ([]byte, error) {
	return nil, nil
}

func (r GetBaseAPI[Req, T]) GetResponse() any {
	if r.Response == nil {
		return new(dto.PtrResponse[T])
	}
	return r.Response
}

func (r GetBaseAPI[Req, T]) GetAction() string {
	return r.Request.GetPath()
}

func (r GetBaseAPI[Req, T]) Decode(bs []byte) error {
	err := json.Unmarshal(bs, r.Response)
	return err
}
