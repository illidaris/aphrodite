package file

import (
	"io"

	"github.com/xuri/excelize/v2"
)

// ExcelExporter 创建一个Exporter，将数据导出到Excel文件。
// sheet: 指定要写入的Sheet名称。
// opts: 一个可选的excelize.Options参数列表，用于配置Excel文件的创建选项。
func ExcelExporter(sheet string, opts ...excelize.Options) Exporter {
	sheet = GetSheetOrDefault(sheet)
	return func(tar io.Writer, headers [][]string, rows ...[]string) error {
		f := excelize.NewFile(opts...)
		defer f.Close()
		allrows := append(headers, rows...)
		for rowIndex, row := range allrows {
			for colIndex, col := range row {
				colId := colIndex + 1
				rowId := rowIndex + 1
				cell, err := excelize.CoordinatesToCellName(colId, rowId)
				if err != nil {
					return err
				}
				err = f.SetCellValue(sheet, cell, col)
				if err != nil {
					return err
				}
			}
		}
		return f.Write(tar)
	}
}

// ExcelImporter 创建一个Importer，从Excel文件中导入数据。
// sheet: 指定要读取的Sheet名称。
// opts: 一个可选的excelize.Options参数列表，用于配置Excel文件的打开选项。
func ExcelImporter(sheet string, opts ...excelize.Options) Importer {
	sheet = GetSheetOrDefault(sheet)
	return func(src io.Reader) ([][]string, error) {
		result := [][]string{}
		f, err := excelize.OpenReader(src)
		if err != nil {
			return result, err
		}
		defer f.Close()
		return f.GetRows(sheet)
	}
}

// GetSheetOrDefault 返回指定的Sheet名称，如果为空，则返回默认Sheet名称。
func GetSheetOrDefault(sheet string) string {
	if sheet == "" {
		return DEF_SHEET_NAME
	}
	return sheet
}
