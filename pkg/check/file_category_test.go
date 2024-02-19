package check

import (
	"io"
	"strings"
	"testing"
)

func TestFileCategory(t *testing.T) {
	tests := []struct {
		filepath string
		category FileCategory
		expected bool
	}{
		{"image.jpg", FileCategoryImage, false},
		{"video.mp4", FileCategoryVideo, false},
		{"audio.wav", FileCategoryAudio, false},
		{"wasm.wasm", FileCategoryWasm, false},
		{"text.txt", FileCategoryNil, false},
		{"", FileCategoryNil, false},
	}

	for _, test := range tests {
		result := FileCategory(test.category).Ok(test.filepath)
		if result != test.expected {
			t.Errorf("Expected %v, but got %v for file %s and category %d", test.expected, result, test.filepath, test.category)
		}
	}
}

func TestFileCategoryRead(t *testing.T) {
	tests := []struct {
		r        io.Reader
		category FileCategory
		expected bool
	}{
		{strings.NewReader("image data"), FileCategoryImage, false},
	}

	for _, test := range tests {
		result := FileCategory(test.category).OkRead(test.r)
		if result != test.expected {
			t.Errorf("Expected %v, but got %v for reader %v and category %d", test.expected, result, test.r, test.category)
		}
	}
}
