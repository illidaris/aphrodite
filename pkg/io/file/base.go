package file

import "errors"

const (
	DEF_SHEET_NAME = "Sheet1"
)

var (
	ErrHeadersNil = errors.New("headers is empty")
	ErrSrcType    = errors.New("src is not string or io.Reader")
)
