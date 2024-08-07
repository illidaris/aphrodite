package file

import (
	"io"
	"os"
)

type Importer func(src io.Reader) ([][]string, error)

func (imr Importer) FromFilePath(src string) ([][]string, error) {
	f, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return imr(f)
}
