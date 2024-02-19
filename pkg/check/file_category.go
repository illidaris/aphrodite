package check

import (
	"io"
	"strings"
)

/*
  if u want more type check, u can use https://github.com/gabriel-vasile/mimetype
*/

// FileCategory represents the category of a file.
type FileCategory int32

// Constants representing the different file categories.
const (
	FileCategoryNil FileCategory = iota
	FileCategoryImage
	FileCategoryVideo
	FileCategoryAudio
	FileCategoryWasm
)

// The Ok function determines whether the given file path belongs to the specified file category.
// Parameters:
//   - filepath: File path
//
// Return value:
//   - bool: Returns true if the file path is of the specified file category; otherwise, returns false.
func (t FileCategory) Ok(filepath string) bool {
	filekeys := strings.Split(filepath, ".")
	if len(filekeys) == 0 {
		return false
	}
	ext := filekeys[len(filekeys)-1]
	if tp := GetFileCategoryByExt(ext); tp != t {
		return false
	}
	f := FileTypeFunc(FileTypeByDetectContentType)
	det, err := f.FrmFile(filepath)
	if err != nil {
		return false
	}
	return GetFileCategoryByDet(det) == t
}

// The OkRead function determines whether the content from the provided io.Reader belongs to the specified file category.
// Parameters:
//   - r: io.Reader
//
// Return value:
//   - bool: Returns true if the content from the io.Reader is of the specified file category; otherwise, returns false.
func (t FileCategory) OkRead(r io.Reader) bool {
	f := FileTypeFunc(FileTypeByDetectContentType)
	det, err := f.FrmReader(r)
	if err != nil {
		return false
	}
	return GetFileCategoryByDet(det) == t
}

// OkBs checks if the given byte slice belongs to the specified file category.
// It first detects the content type of the byte slice using FileTypeByDetectContentType.
// Then it checks if the detected file category matches the specified file category.
// Returns true if the byte slice belongs to the specified file category, otherwise false.
func (t FileCategory) OkBs(r []byte) bool {
	det := FileTypeByDetectContentType(r)
	return GetFileCategoryByDet(det) == t
}
