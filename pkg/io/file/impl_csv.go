package file

import (
	"encoding/csv"
	"io"

	"golang.org/x/text/transform"
)

const (
	CSV_BOM = "\xEF\xBB\xBF" // 写入UTF-8 BOM，避免使用Microsoft Excel打开乱码
)

// CsvExporter 创建并返回一个Exporter类型的函数，用于将数据写入CSV格式。
// 该函数接受一个io.Writer类型的参数用于指定输出流，接受一个二维字符串数组headers和一个可变长度的字符串数组rows作为数据。
// 它首先写入headers中的每一行，然后写入rows中的每一行到CSV文件中。
// 最后，它确保CSV写入器被正确地刷新，以保证所有数据都被写入。
func CsvExporter() Exporter {
	return func(tar io.Writer, headers [][]string, rows ...[]string) error {
		wt := csv.NewWriter(tar)
		for _, header := range headers {
			_ = wt.Write(header)
		}
		for _, header := range rows {
			_ = wt.Write(header)
		}
		wt.Flush()
		return nil
	}
}

// CsvImporter 创建并返回一个Importer类型的函数，用于从CSV文件中导入数据。
// 它接受一个io.Reader类型的参数src和一个可选的transform.Transformer类型的参数t。
// 如果提供了转换器t，则使用transform.NewReader对src进行转换。
// 然后，它使用csv.NewReader从src中读取所有数据，并以二维字符串数组的形式返回这些数据。
func CsvImporter(t transform.Transformer) Importer {
	return func(src io.Reader) ([][]string, error) {
		if t != nil {
			src = transform.NewReader(src, t)
		}
		reader := csv.NewReader(src)
		return reader.ReadAll()
	}
}
