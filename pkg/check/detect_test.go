package check

import "testing"

func TestFileTypeByDetectContentType(t *testing.T) {
	bs := []byte("Hello, World! TestFileTypeByDetectContentType")
	expected := "text/plain; charset=utf-8"
	result := FileTypeByDetectContentType(bs)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
