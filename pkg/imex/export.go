package imex

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/illidaris/aphrodite/pkg/convert/table2struct"
	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/illidaris/aphrodite/pkg/group"
	"github.com/xuri/excelize/v2"
)

// BaseExport 导出数据
// 参数：
//
//	ctx context.Context 上下文
//	req Req 初始请求
//	exportFunc func(context.Context, Req) (Resp, error) 获取导出数据函数
//	pagesFunc func(Req, Resp) []Req 获取分页数据函数，利用初始请求，产出所有分页请求
//	getItemFunc func(Resp) []any 解析【获取导出数据函数】的响应，转化出具体数据结构
//	opts ...ImExOptionFunc[Resp] 选项
func BaseExport[Req dependency.IPage, Resp dependency.IPaginator, Item any](
	ctx context.Context,
	req Req,
	exportFunc func(context.Context, Req) (Resp, error),
	pagesFunc func(Req, Resp) []Req,
	getItemFunc func(Resp) []Item,
	opts ...ImExOptionFunc[Item],
) (io.Reader, string, error) {
	opt := NewImExOption[Item]()
	for _, f := range opts {
		f(opt)
	}
	firstResp, ex := exportFunc(ctx, req)
	if ex != nil {
		return nil, opt.ExportName, ex
	}
	pages := pagesFunc(req, firstResp)
	inCh := make(chan Item, 10)
	go func() {
		defer close(inCh)
		for _, v := range getItemFunc(firstResp) {
			inCh <- v
		}
		_, _ = group.GroupBaseFunc(func(subReqs ...Req) (int64, error) {
			affect := 0
			for _, subReq := range subReqs {
				resp, ex := exportFunc(ctx, subReq)
				if ex != nil {
					continue
				}
				for _, v := range getItemFunc(resp) {
					inCh <- v
				}
				affect++
			}
			return int64(affect), nil
		}, 1, pages...)
	}()
	return BaseExportToReader(ctx, inCh, opt)
}

func BaseExportToReader[T any](ctx context.Context, inCh <-chan T, opt *ImExOption[T]) (io.Reader, string, error) {
	items := []any{}
	for v := range inCh {
		items = append(items, v)
	}
	reader, err := WriteReader(ctx, items, opt)
	return reader, opt.ExportName, err
}

func WriteReader[T any](ctx context.Context, items []any, opt *ImExOption[T]) (io.Reader, error) {
	allRows := [][]string{}
	headers, rows, err := table2struct.Struct2Table(items, opt.Table2StructOptions...)
	if err != nil {
		return nil, err
	}
	allRows = append(allRows, headers...)
	allRows = append(allRows, rows...)

	nameKeys := strings.Split(opt.ExportName, ".")
	if len(nameKeys) < 2 {
		return nil, errors.New("文件格式错误")
	}
	ext := nameKeys[len(nameKeys)-1]
	switch ext {
	case "xlsx", "xls":
		f := excelize.NewFile()
		for rowIndex, row := range allRows {
			for colIndex, v := range row {
				cellName, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
				_ = f.SetCellValue("Sheet1", cellName, v)
			}
		}
		return f.WriteToBuffer()
	case "csv":
		bs := []byte{}
		w := bytes.NewBuffer(bs)
		_, _ = w.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，防止中文乱码
		csvW := csv.NewWriter(w)
		csvW.UseCRLF = true
		if err := csvW.WriteAll(rows); err != nil {
			return w, err
		}
		csvW.Flush()
		return w, nil
	case "json":
		bs, err := json.Marshal(items)
		if err != nil {
			return bytes.NewBuffer(bs), err
		}
		return bytes.NewBuffer(bs), nil
	}
	return nil, errors.New("文件格式错误")
}
