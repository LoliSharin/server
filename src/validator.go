package main

import (
	"errors"
	"path/filepath"
	"strings"
)

var (
	FileSizeLimit    = errors.New("размер файла не допустимый")
	FileTypeError    = errors.New("не подходящий формат файла")
	FileNameNotEmpty = errors.New("файл не должен быть пустым")
)

const MaxFileSize = 50 * 1024 * 1024 // 50 MB
const maxFileSize = MaxFileSize

func ValidateMedia(filename string, size int64) error {
	if strings.TrimSpace(filename) == "" {
		return FileNameNotEmpty
	}
	if size <= 0 || size > MaxFileSize {
		return FileSizeLimit
	}
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".mp4", ".mov", ".webm", ".png", ".jpg", ".jpeg":
		return nil
	default:
		return FileTypeError
	}
}
