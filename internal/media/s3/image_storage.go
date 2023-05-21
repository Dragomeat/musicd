package s3

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"musicd/internal/errors"
	mediaImage "musicd/internal/media/image"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
	_ "golang.org/x/image/webp"
)

type ErrUnknownContentType struct {
	ContentType string
}

func (e ErrUnknownContentType) Code() errors.ErrorCode {
	return "media.unknownContentType"
}

func (e ErrUnknownContentType) Message() string {
	return fmt.Sprintf("unknown content type %s", e.ContentType)
}

func (e ErrUnknownContentType) Error() string {
	return fmt.Sprintf("unknown content type %s", e.ContentType)
}

var contentTypeToExtension = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpeg",
	"image/webp": ".webp",
}

const (
	mb = 1 << 20
)

type ImageStorage struct {
	client         *s3.Client
	urlConstructor *ImageUrlConstructor
	logger         *zap.Logger
	bucket         string
}

func NewImageStorage(
	client *s3.Client,
	urlConstructor *ImageUrlConstructor,
	logger *zap.Logger,
	bucket string,
) *ImageStorage {
	return &ImageStorage{
		client:         client,
		bucket:         bucket,
		urlConstructor: urlConstructor,
		logger:         logger,
	}
}

func (u *ImageStorage) Upload(
	ctx context.Context,
	upload mediaImage.Upload,
) (mediaImage.Image, error) {
	img := mediaImage.Image{
		FileName: upload.FileName,
		Size:     upload.Size,
	}

	metadata, reader, err := u.parseMetadata(upload.File)
	if err != nil {
		return img, err
	}

	if _, ok := contentTypeToExtension[metadata.ContentType]; !ok {
		return img, ErrUnknownContentType{ContentType: metadata.ContentType}
	}

	img.Sizes = metadata.Sizes
	img.ContentType = metadata.ContentType

	// Here we most likely have already validated file for the max size, content type, etc ...
	// And we need whole file in memory to generate id and upload this file
	content, err := io.ReadAll(reader)
	if err != nil {
		return img, err
	}

	img.Id, err = u.generateId(content, img.ContentType)
	if err != nil {
		return img, err
	}

	err = u.upload(ctx, bytes.NewBuffer(content), img)
	return img, err
}

type imageMetadata struct {
	Sizes       mediaImage.Sizes
	ContentType string
}

func (u *ImageStorage) parseMetadata(content io.Reader) (imageMetadata, io.Reader, error) {
	metadata := imageMetadata{}

	header := new(bytes.Buffer)
	teeHeader := io.TeeReader(io.LimitReader(content, mb), header)
	content = io.MultiReader(header, content)

	config, _, err := image.DecodeConfig(teeHeader)
	if err != nil {
		return metadata, content, err
	}

	sizes := mediaImage.Sizes{
		Width:  config.Width,
		Height: config.Height,
	}
	metadata.Sizes = sizes
	metadata.ContentType = http.DetectContentType(header.Bytes())

	return metadata, content, err
}

func (u *ImageStorage) generateId(content []byte, contentType string) (string, error) {
	ext, ok := contentTypeToExtension[contentType]
	if !ok {
		return "", ErrUnknownContentType{ContentType: contentType}
	}

	hasher := sha1.New()
	_, err := hasher.Write(content)
	if err != nil {
		return "", err
	}

	hash := hex.EncodeToString(hasher.Sum(nil))
	return string(hash[0]) + "/" + string(hash[1]) + "/" + hash[2:] + ext, nil
}

func (u *ImageStorage) upload(ctx context.Context, reader io.Reader, image mediaImage.Image) error {
	_, err := u.client.PutObject(
		ctx, &s3.PutObjectInput{
			Bucket:      aws.String(u.bucket),
			Key:         aws.String(image.Id),
			Body:        reader,
			ContentType: aws.String(image.ContentType),
			Metadata: map[string]string{
				"dimensions": fmt.Sprintf("%dx%d", image.Sizes.Width, image.Sizes.Height),
			},
		},
	)

	return err
}
