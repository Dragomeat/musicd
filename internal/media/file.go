package media

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type metadata struct {
	OriginalName string `json:"originalName"`
	ContentType  string `json:"contentType"`
	Size         int    `json:"size"`
}

type FileDriver struct {
	basePath string
}

func NewFileDriver(basePath string) *FileDriver {
	return &FileDriver{basePath: basePath}
}

func (d *FileDriver) GetMetadata(ctx context.Context, id string) (File, error) {
	mf, err := os.Open(d.pathToMetadata(id))
	if err != nil {
		return File{}, err
	}
	defer mf.Close()

	var md metadata
	err = json.NewDecoder(mf).Decode(&md)
	if err != nil {
		return File{}, err
	}

	return File{
		Id:          id,
		FileName:    md.OriginalName,
		ContentType: md.ContentType,
		Size:        md.Size,
	}, nil
}

func (d *FileDriver) Put(ctx context.Context, file File, content io.Reader) error {
	mf, err := os.Create(d.pathToMetadata(file.Id))
	if err != nil {
		return err
	}
	defer mf.Close()

	md := metadata{
		OriginalName: file.FileName,
		ContentType:  file.ContentType,
		Size:         file.Size,
	}

	err = json.NewEncoder(mf).Encode(md)
	if err != nil {
		return err
	}

	cf, err := os.Create(d.pathToContent(file.Id))
	if err != nil {
		return err
	}
	defer cf.Close()

	_, err = io.Copy(cf, content)
	if err != nil {
		return err
	}

	return nil
}

func (d *FileDriver) pathToContent(id string) string {
	return fmt.Sprintf("%s/%s_ct", d.basePath, id)
}

func (d *FileDriver) pathToMetadata(id string) string {
	return fmt.Sprintf("%s/%s_md", d.basePath, id)
}
