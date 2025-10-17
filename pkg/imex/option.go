package imex

import "github.com/illidaris/aphrodite/pkg/convert/table2struct"

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
	Deep                bool
}

func (o ImExOption[T]) Table2Struct(dst interface{}, rows [][]string) (err error) {
	if !o.Deep {
		return table2struct.Table2Struct(dst, rows, o.Table2StructOptions...)
	}
	return table2struct.Table2Objs(dst, rows, o.Table2StructOptions...)
}

func (o ImExOption[T]) Struct2Table(dsts []interface{}) ([][]string, [][]string, error) {
	if !o.Deep {
		return table2struct.Struct2Table(dsts, o.Table2StructOptions...)
	}
	return table2struct.Objs2Table(dsts, o.Table2StructOptions...)
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

func WithDeep[T any]() ImExOptionFunc[T] {
	return func(opt *ImExOption[T]) {
		opt.Deep = true
		if opt.Table2StructOptions == nil {
			opt.Table2StructOptions = []table2struct.Table2StructOptionFunc{}
		}
		opt.Table2StructOptions = append(opt.Table2StructOptions, table2struct.WithDeep())
	}
}
