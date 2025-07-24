package api

import "net/url"

type IRequest interface {
	GetAction() string
}

type IGetUrlQuery interface {
	GetUrlQuery() url.Values
}
