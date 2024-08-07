package file

import (
	"encoding/json"
	"io"
)

func JsonExporter(pretty bool) Exporter {
	return func(tar io.Writer, headers [][]string, rows ...[]string) error {
		if len(headers) == 0 {
			return ErrHeadersNil
		}
		headerMap := map[int]string{}
		for k, v := range headers[0] {
			headerMap[k] = v
		}
		result := make([]map[string]interface{}, 0)
		for _, row := range rows {
			m := make(map[string]interface{})
			for rowIndex, v := range row {
				key, ok := headerMap[rowIndex]
				if !ok {
					continue
				}
				m[key] = v
			}
			result = append(result, m)
		}
		if pretty {
			bs, err := json.MarshalIndent(result, "", "    ")
			if err != nil {
				return err
			}
			tar.Write(bs)
		} else {
			bs, err := json.Marshal(result)
			if err != nil {
				return err
			}
			tar.Write(bs)
		}
		return nil
	}
}

// TODO: 暂未实现
func JsonImporter() Importer {
	return func(src io.Reader) ([][]string, error) {
		panic("no implementation")
	}
}
