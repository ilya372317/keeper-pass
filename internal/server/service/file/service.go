package file

import (
	"context"
	"fmt"
	"io"
)

type fileStorage interface {
	SaveFile(ctx context.Context, filePath string, file io.Reader) error
}

type Service struct {
	fileStorage fileStorage
}

func New(fileStorage fileStorage) *Service {
	return &Service{fileStorage: fileStorage}
}

func (s *Service) Upload(ctx context.Context, filePath string, file io.Reader) error {
	if err := s.fileStorage.SaveFile(ctx, filePath, file); err != nil {
		return fmt.Errorf("failed upload file: %w", err)
	}

	return nil
}
