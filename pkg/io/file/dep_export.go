package file

import (
	"io"
	"os"
)

type Exporter func(tar io.Writer, headers [][]string, rows ...[]string) error

func (exr Exporter) ToFilePath(target string, headers [][]string, rows ...[]string) error {
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()
	return exr(f, headers, rows...)
}
