package media

import (
	"context"
	"errors"
	"fmt"
	"io"
)

var ErrFileNotFound = errors.New("file not found")

type Upload struct {
	Id          string
	File        io.Reader
	FileName    string
	ContentType string
	Size        int
}

type File struct {
	Id          string `json:"id"`
	Sha256      string `json:"sha256"`
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	Size        int    `json:"size"`
}

//go:generate go run github.com/vektra/mockery/v2 --name Uploader
type Uploader interface {
	Upload(ctx context.Context, upload Upload) (File, error)
}

//go:generate go run github.com/vektra/mockery/v2 --name Driver
type Driver interface {
	GetMetadata(ctx context.Context, id string) (File, error)
	Put(ctx context.Context, upload Upload) (File, error)
	Move(ctx context.Context, oldKey string, newKey string) error
}

type Storage struct {
	driver Driver
}

func NewStorage(driver Driver) *Storage {
	return &Storage{driver: driver}
}

func (s *Storage) Upload(ctx context.Context, upload Upload) (File, error) {
	file, err := s.driver.Put(ctx, upload)
	if err != nil {
		return File{}, err
	}

	err = s.driver.Move(ctx, upload.Id, file.Sha256)
	if err != nil {
		return File{}, err
	}

	file.Id = file.Sha256

	return file, nil
}

func (s *Storage) GetLink(ctx context.Context, id string) (string, error) {
	file, err := s.driver.GetMetadata(ctx, id)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://cdn.music.local/%s.%s", file.Id, "mp3"), nil
}
