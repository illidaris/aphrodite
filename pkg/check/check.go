package check

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// SaveUploadedFileWithCheck saves an uploaded file while performing checks.
func SaveUploadedFileWithCheck(file *multipart.FileHeader, dst string, allowCategorys ...FileCategory) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Read the entire file content
	bs, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	// Extract the file extension
	filekeys := strings.Split(file.Filename, ".")
	if len(filekeys) == 0 {
		return errors.New("no file name found")
	}
	ext := filekeys[len(filekeys)-1]

	// Determine the file category based on its extension
	tp := GetFileCategoryByExt(ext)

	// Check if the file type is allowed
	isAllowed := false
	for _, cgy := range allowCategorys {
		if tp == cgy {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return errors.New("file extension is not allowed")
	}

	// Validate the file content against its type
	if !tp.OkBs(bs) {
		return errors.New("file type mismatch")
	}

	// Create directories for the destination path if needed
	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	// Write the file to the destination
	return os.WriteFile(dst, bs, os.ModePerm)
}
