package api

import (
	"encoding/json"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/illidaris/aphrodite/dto"
	restSender "github.com/illidaris/rest/sender"
)

func NewPostAPI[Req IRequest, T any](param Req) *PostBaseAPI[Req, T] {
	return &PostBaseAPI[Req, T]{
		Request:  param,
		Response: new(dto.PtrResponse[T]),
	}
}

var _ = restSender.IRequest(&PostBaseAPI[IRequest, any]{})

type PostBaseAPI[Req IRequest, T any] struct {
	restSender.JSONRequest `json:"-"`
	Request                Req
	Response               *dto.PtrResponse[T] `json:"-"`
}

func (r PostBaseAPI[Req, T]) GetUrlQuery() url.Values {
	u, _ := query.Values(r.Request)
	return u
}

func (r PostBaseAPI[Req, T]) Encode() ([]byte, error) {
	return json.Marshal(r.Request)
}

func (r PostBaseAPI[Req, T]) GetResponse() any {
	if r.Response == nil {
		return new(dto.PtrResponse[T])
	}
	return r.Response
}

func (r PostBaseAPI[Req, T]) GetAction() string {
	return r.Request.GetAction()
}

func (r PostBaseAPI[Req, T]) Decode(bs []byte) error {
	err := json.Unmarshal(bs, r.Response)
	return err
}
