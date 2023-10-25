package exception

import (
	"fmt"
	"strings"
)

type Exception interface {
	Code() int32
	error
}

func New(ex ExceptionType, s string) Exception {
	return &errorString{
		ex: ex,
		s:  s,
	}
}

func Wrap(ex ExceptionType, err error, msgs ...string) Exception {
	switch len(msgs) {
	case 0:
		return &errorString{
			ex:  ex,
			err: err,
		}
	case 1:
		return &errorString{
			ex:  ex,
			err: err,
			s:   msgs[0],
		}
	default:
		return &errorString{
			ex:  ex,
			err: err,
			s:   strings.Join(msgs, ","),
		}
	}
}

var _ = error(&errorString{})
var _ = Exception(&errorString{})

type errorString struct {
	ex  ExceptionType
	err error
	s   string
}

func (e *errorString) Code() int32 {
	return int32(e.ex)
}

func (e *errorString) Error() string {
	if e.err != nil {
		if e.s == "" {
			return e.err.Error()
		}
		return fmt.Sprintf("%s,%s", e.s, e.err.Error())
	}
	return e.s
}
