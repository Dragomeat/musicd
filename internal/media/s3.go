package media

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type writerFunc func(p []byte) (n int, err error)

func (f writerFunc) Write(p []byte) (n int, err error) {
	return f(p)
}

type S3Driver struct {
	s3     *s3.Client
	bucket string
}

func NewS3Driver(s3 *s3.Client, bucket string) *S3Driver {
	return &S3Driver{s3: s3, bucket: bucket}
}

func (d *S3Driver) GetMetadata(ctx context.Context, id string) (File, error) {
	return File{}, nil
}

func (d *S3Driver) Put(ctx context.Context, upload Upload) (File, error) {
	uploader := manager.NewUploader(d.s3)
	hasher := sha256.New()
	size := 0
	res, err := uploader.Upload(
		ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(d.bucket),
			Key:    aws.String(upload.Id),
			Body: io.TeeReader(
				upload.File,
				io.MultiWriter(
					hasher,
					writerFunc(func(p []byte) (n int, err error) {
						n = len(p)
						size += n
						if size > upload.Size {
							return 0, fmt.Errorf("expected size %d, got %d", upload.Size, size)
						}
						return n, nil
					},
					),
				),
			),
			ContentType: aws.String(upload.ContentType),
			Metadata: map[string]string{
				"originalName": upload.FileName,
			},
		},
	)
	if err != nil {
		return File{}, err
	}

	hash := hex.EncodeToString(hasher.Sum(nil))
	if res.ChecksumSHA256 != nil {
		return File{}, fmt.Errorf("expected checksum %s, got %s", hash, *res.ChecksumSHA256)
	}

	return File{
		Id:          upload.Id,
		Sha256:      hash,
		FileName:    upload.FileName,
		ContentType: upload.ContentType,
		Size:        upload.Size,
	}, err
}

func (d *S3Driver) Move(ctx context.Context, oldKey string, newKey string) error {
	_, err := d.s3.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(d.bucket),
		CopySource: aws.String(fmt.Sprintf("%s/%s", d.bucket, oldKey)),
		Key:        aws.String(newKey),
	})
	if err != nil {
		return err
	}

	_, err = d.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(oldKey),
	})
	return err
}
