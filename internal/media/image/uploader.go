package image

import (
	"context"
	"io"
)

type Upload struct {
	File     io.Reader
	FileName string
	Size     int64
}

type Uploader interface {
	Upload(ctx context.Context, upload Upload) (Image, error)
}
