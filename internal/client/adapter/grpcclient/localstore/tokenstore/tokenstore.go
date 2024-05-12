package tokenstore

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/ilya372317/pass-keeper/internal/client/domain"
)

type TokenStorage struct {
	credFile *os.File
}

func New(credFile *os.File) *TokenStorage {
	return &TokenStorage{
		credFile: credFile,
	}
}

func (u *TokenStorage) SetAccessToken(token string) error {
	fileStat, err := u.credFile.Stat()
	if err != nil {
		return fmt.Errorf("failed get stat of file with access token: %w", err)
	}
	fileStat.Size()
	if err = u.credFile.Truncate(fileStat.Size()); err != nil {
		return fmt.Errorf("failed truncate old token in file storage: %w", err)
	}
	if _, err = u.credFile.Write(append([]byte(token), '\n')); err != nil {
		return fmt.Errorf("failed write new token to file storage: %w", err)
	}

	return nil
}

func (u *TokenStorage) GetAccessToken() (string, error) {
	if _, err := u.credFile.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("failed rewind access token file storage: %w", err)
	}
	buff := bufio.NewReader(u.credFile)

	token, err := buff.ReadBytes('\n')
	if err != nil {
		return "", domain.ErrUnauthenticated
	}

	return string(token[:len(token)-1]), nil
}
