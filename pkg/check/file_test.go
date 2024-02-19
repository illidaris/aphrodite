package check

import (
	"bytes"
	"os"
	"testing"
)

func TestFileTypeFrmReader(t *testing.T) {
	fileType := FileTypeFunc(FileTypeByDetectContentType)
	reader := bytes.NewReader([]byte("Hello, World! TestFileTypeFrmReader"))
	expected := "text/plain; charset=utf-8"
	result, err := fileType.FrmReader(reader)
	if err != nil {
		t.Error(err)
	}
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestFileTypeFrmFile(t *testing.T) {
	fileType := FileTypeFunc(FileTypeByDetectContentType)
	file, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	file.WriteString("Hello, World! TestFileTypeFrmFile")
	expected := "text/plain; charset=utf-8"
	result, err := fileType.FrmFile(file.Name())
	if err != nil {
		t.Error(err)
	}
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
