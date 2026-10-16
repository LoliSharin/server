package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"time"
)

type StorageService struct {
	repo      *FileRepo
	uploadDir string
}

func NewStorageService(repo *FileRepo, uploadDir string) *StorageService {
	return &StorageService{
		repo:      repo,
		uploadDir: uploadDir,
	}
}

func (s *StorageService) GetFileForStream(ctx context.Context, id int) (*os.File, *FileMetadata, error) {
	meta, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	readFile, err := os.Open(meta.FilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("файл не найден на диске: %w", err)
	}
	return readFile, meta, nil
}

func (s *StorageService) SaveUploadedFile(
	ctx context.Context,
	fileName string,
	mimeType string,
	src io.Reader,
) (*FileMetadata, error) {
	safeName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(fileName))
	fullPath := filepath.Join(s.uploadDir, safeName)

	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать файл на диске: %w", err)
	}

	// Флаг успешного сохранения: при любой ошибке файл закрывается и удаляется с диска
	var success bool
	defer func() {
		if !success {
			_ = dst.Close()
			_ = os.Remove(fullPath)
		}
	}()

	hasher := sha256.New()
	tee := io.TeeReader(src, hasher)

	size, err := io.Copy(dst, tee)
	if err != nil {
		return nil, fmt.Errorf("ошибка записи потока данных: %w", err)
	}

	if err := dst.Close(); err != nil {
		return nil, fmt.Errorf("ошибка закрытия файла на диске: %w", err)
	}

	if mimeType == "" {
		mimeType = mime.TypeByExtension(filepath.Ext(fileName))
	}

	hashSum := hex.EncodeToString(hasher.Sum(nil))

	meta := &FileMetadata{
		FileName: fileName,
		FilePath: fullPath,
		Size:     size,
		MIMEType: mimeType,
		SHA256:   hashSum,
	}

	if err := s.repo.Save(ctx, meta); err != nil {
		return nil, fmt.Errorf("не удалось сохранить метаданные: %w", err)
	}

	success = true
	return meta, nil
}
