package exception

import (
	"fmt"
	"strings"
)

type Exception interface {
	Code() int32
	error
}

// New函数用于创建一个Exception对象，该对象包含一个异常类型ex和一个字符串s。
func New(ex ExceptionType, s string) Exception {
	return &errorString{
		ex: ex,
		s:  s,
	}
}

// Wrap函数用于创建一个Exception类型的对象，该对象包含了一个异常类型ex、一个错误err和一个可变参数msgs。
// 当msgs参数为空时，返回一个errorString对象，该对象的ex字段为ex，err字段为err。
// 当msgs参数只有一个元素时，返回一个errorString对象，该对象的ex字段为ex，err字段为err，s字段为msgs[0]。
// 当msgs参数有多个元素时，返回一个errorString对象，该对象的ex字段为ex，err字段为err，s字段为msgs所有元素用逗号连接起来的字符串。
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
