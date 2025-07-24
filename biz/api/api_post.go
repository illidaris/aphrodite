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
		Response: new(dto.TResponse[T]),
	}
}

var _ = restSender.IRequest(&PostBaseAPI[IRequest, any]{})

type PostBaseAPI[Req IRequest, T any] struct {
	restSender.JSONRequest `json:"-"`
	Request                Req
	Response               *dto.TResponse[T] `json:"-"`
}

func (r PostBaseAPI[Req, T]) GetUrlQuery() url.Values {
	v, ok := any(r.Request).(IGetUrlQuery)
	if ok {
		return v.GetUrlQuery()
	}
	u, _ := query.Values(r.Request)
	return u
}

func (r PostBaseAPI[Req, T]) Encode() ([]byte, error) {
	return json.Marshal(r.Request)
}

func (r PostBaseAPI[Req, T]) GetResponse() any {
	if r.Response == nil {
		return new(dto.TResponse[T])
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

func NewFormAPI[Req IRequest, T any](param Req) *PostBaseAPI[Req, T] {
	return &PostBaseAPI[Req, T]{
		Request:  param,
		Response: new(dto.TResponse[T]),
	}
}

var _ = restSender.IRequest(&PostBaseAPI[IRequest, any]{})

type FormBaseAPI[Req IRequest, T any] struct {
	restSender.FormUrlEncodeRequest `json:"-"`
	Request                         Req
	Response                        *dto.PtrResponse[T] `json:"-"`
	Queries                         url.Values          `json:"-"`
}

func (r FormBaseAPI[Req, T]) GetUrlQuery() url.Values {
	return r.Queries
}

func (r FormBaseAPI[Req, T]) Encode() ([]byte, error) {
	u, err := query.Values(r.Request)
	if err != nil {
		return nil, err
	}
	return []byte(u.Encode()), nil
}

func (r FormBaseAPI[Req, T]) GetResponse() any {
	if r.Response == nil {
		return new(dto.TResponse[T])
	}
	return r.Response
}

func (r FormBaseAPI[Req, T]) GetAction() string {
	return r.Request.GetAction()
}

func (r FormBaseAPI[Req, T]) Decode(bs []byte) error {
	err := json.Unmarshal(bs, r.Response)
	return err
}
