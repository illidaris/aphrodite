package check

import (
	"fmt"
	"io"
	"os"
)

// FileTypeFunc is a function type that takes a byte slice and returns a string.
type FileTypeFunc func(bs []byte) string

// FrmReader reads from the given reader and returns the file type using the FileTypeFunc.
func (f FileTypeFunc) FrmReader(r io.Reader) (string, error) {
	bs := make([]byte, 512)
	n, err := r.Read(bs)
	if err != nil || n < 1 {
		return "", fmt.Errorf("failed to read from reader: %w", err)
	}
	if n < 512 {
		bs = bs[:n]
	}
	return f(bs), nil
}

// FrmFile reads from the given file and returns the file type using the FileTypeFunc.
func (f FileTypeFunc) FrmFile(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	return f.FrmReader(file)
}
