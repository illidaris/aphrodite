package imex

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/illidaris/aphrodite/pkg/convert/table2struct"
	"github.com/xuri/excelize/v2"
)

type IImpoert interface {
	GetReader() io.Reader
	GetFileName() string
}

func NewImExOption[T any]() *ImExOption[T] {
	return &ImExOption[T]{
		Table2StructOptions: make([]table2struct.Table2StructOptionFunc, 0),
		Iterates:            make([]func(item *T), 0),
	}
}

type ImExOption[T any] struct {
	Table2StructOptions []table2struct.Table2StructOptionFunc
	Iterates            []func(item *T)
	ExportName          string
}
type ImExOptionFunc[T any] func(opt *ImExOption[T])

func WithTable2StructOptions[T any](fs ...table2struct.Table2StructOptionFunc) ImExOptionFunc[T] {
	return func(opt *ImExOption[T]) {
		if opt.Table2StructOptions == nil {
			opt.Table2StructOptions = make([]table2struct.Table2StructOptionFunc, 0)
		}
		opt.Table2StructOptions = append(opt.Table2StructOptions, fs...)
	}
}

func WithIterates[T any](fs ...func(*T)) ImExOptionFunc[T] {
	return func(opt *ImExOption[T]) {
		if opt.Iterates == nil {
			opt.Iterates = make([]func(item *T), 0)
		}
		opt.Iterates = append(opt.Iterates, fs...)
	}
}

func WithExportName[T any](name string) ImExOptionFunc[T] {
	return func(opt *ImExOption[T]) {
		opt.ExportName = name
	}
}
func BaseImport[Req IImpoert, Resp any](ctx context.Context, req Req, opts ...ImExOptionFunc[Resp]) ([]*Resp, error) {
	opt := NewImExOption[Resp]()
	for _, f := range opts {
		f(opt)
	}
	awards, err := ParseReader(ctx, req.GetReader(), req.GetFileName(), opt)
	if err != nil {
		return nil, err
	}
	for _, iterate := range opt.Iterates {
		for _, v := range awards {
			iterate(v)
		}
	}
	return awards, nil
}

func ParseReader[T any](ctx context.Context, reader io.Reader, filename string, opt *ImExOption[T]) ([]*T, error) {
	result := []*T{}
	nameKeys := strings.Split(filename, ".")
	if len(nameKeys) < 2 {
		return result, errors.New("文件格式错误")
	}
	ext := nameKeys[len(nameKeys)-1]
	switch ext {
	case "xlsx", "xls":
		excel, err := excelize.OpenReader(reader)
		if err != nil {
			return result, err
		}
		rows, err := excel.GetRows("Sheet1")
		if err != nil {
			return result, err
		}
		err = table2struct.Table2Struct(&result, rows, opt.Table2StructOptions...)
		if err != nil {
			return result, err
		}
		return result, nil
	case "csv":
		csvReader := csv.NewReader(reader)
		rows, err := csvReader.ReadAll()
		if err != nil {
			return result, err
		}
		err = table2struct.Table2Struct(&result, rows, opt.Table2StructOptions...)
		if err != nil {
			return result, err
		}
		return result, nil
	case "json":
		bs, err := io.ReadAll(reader)
		if err != nil {
			return result, err
		}
		err = json.Unmarshal(bs, &result)
		if err != nil {
			return result, err
		}
		return result, nil
	}
	return result, errors.New("文件格式错误")
}
