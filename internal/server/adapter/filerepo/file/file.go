package file

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

const BucketName = "pass"

type Storage struct {
	c *minio.Client
}

func New(c *minio.Client) *Storage {
	return &Storage{
		c: c,
	}
}

func (s *Storage) SaveFile(ctx context.Context, filePath string, file io.Reader) error {
	if _, err := s.c.PutObject(ctx, BucketName, filePath, file, -1, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	}); err != nil {
		return fmt.Errorf("failed upload file to minio: %w", err)
	}

	return nil
}
