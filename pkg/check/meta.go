package check

import "strings"

// GetFileCategoryByDet returns the FileCategory for a given detection string.
func GetFileCategoryByDet(det string) FileCategory {
	return DetCgyMap[det]
}

// DetCgyMap is a map of file extensions to their corresponding FileCategory.
var DetCgyMap = map[string]FileCategory{
	"image/x-icon":     FileCategoryImage,
	"image/bmp":        FileCategoryImage,
	"image/gif":        FileCategoryImage,
	"image/webp":       FileCategoryImage,
	"image/png":        FileCategoryImage,
	"image/apng":       FileCategoryImage,
	"image/jpeg":       FileCategoryImage,
	"video/webm":       FileCategoryVideo,
	"video/mp4":        FileCategoryVideo,
	"video/avi":        FileCategoryVideo,
	"application/ogg":  FileCategoryVideo,
	"audio/mpeg":       FileCategoryAudio,
	"audio/aiff":       FileCategoryAudio,
	"audio/midi":       FileCategoryAudio,
	"audio/wave":       FileCategoryAudio,
	"application/wasm": FileCategoryWasm,
}

// GetFileCategoryByExt returns the FileCategory for a given file extension.
func GetFileCategoryByExt(ext string) FileCategory {
	return ExtCgyMap[strings.ToLower(ext)]
}

// ExtCgyMap is a map of file extensions to their corresponding FileCategory.
var ExtCgyMap = map[string]FileCategory{
	"webp": FileCategoryImage,
	"jpg":  FileCategoryImage,
	"png":  FileCategoryImage,
	"apng": FileCategoryImage,
	"gif":  FileCategoryImage,
	"bmp":  FileCategoryImage,
	"ico":  FileCategoryImage,
	"webm": FileCategoryVideo,
	"mp4":  FileCategoryVideo,
	"avi":  FileCategoryVideo,
	"ogg":  FileCategoryVideo,
	"mpeg": FileCategoryAudio,
	"mp3":  FileCategoryAudio,
	"aiff": FileCategoryAudio,
	"midi": FileCategoryAudio,
	"wav":  FileCategoryAudio,
	"wasm": FileCategoryWasm,
}
