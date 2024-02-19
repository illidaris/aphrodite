package check

import "net/http"

// FileTypeByDetectContentType detects the content type of the given byte slice.
func FileTypeByDetectContentType(bs []byte) string {
	return http.DetectContentType(bs)
}
