package media

import (
	"context"
	"musicd/internal/logger"
	"musicd/internal/media/api"
	"musicd/internal/media/image"
	"musicd/internal/media/s3"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ModuleConfig struct {
	S3Endpoint        string `mapstructure:"S3_ENDPOINT"`
	S3Region          string `mapstructure:"S3_REGION"`
	S3AccessKey       string `mapstructure:"S3_ACCESS_KEY"`
	S3AccessSecret    string `mapstructure:"S3_ACCESS_SECRET"`
	S3BucketForImages string `mapstructure:"S3_BUCKET_FOR_IMAGES"`

	Storage          string `mapstructure:"MEDIA_STORAGE"`
	CdnHost          string `mapstructure:"MEDIA_CDN_HOST"`
	ProtectionSecret string `mapstructure:"MEDIA_PROTECTION_SECRET"`
}

func NewModule() fx.Option {
	return fx.Module(
		"media",
		fx.Provide(
			func(viper *viper.Viper) (*ModuleConfig, error) {
				config := ModuleConfig{}
				err := viper.Unmarshal(&config)
				if err != nil {
					return nil, err
				}
				return &config, nil
			},

			api.NewTransformer,

			NewStorage,
			func() *FileDriver {
				return NewFileDriver("/tmp")
			},
			func(client *awsS3.Client) *S3Driver {
				return NewS3Driver(client, "musicd")
			},
			func(driver *S3Driver) Driver {
				return driver
			},
			func(storage *Storage) Uploader {
				return storage
			},

			func(logger logger.Logger, config *ModuleConfig) (*awsS3.Client, error) {
				cfg := aws.Config{
					Logger: logging.LoggerFunc(func(classification logging.Classification, format string, v ...interface{}) {
						logger.Info(context.Background(), format, v...)
					}),
					Credentials: credentials.NewStaticCredentialsProvider(
						config.S3AccessKey,
						config.S3AccessSecret,
						"",
					),
					Region: config.S3Region,
					EndpointResolverWithOptions: aws.EndpointResolverWithOptionsFunc(
						func(service, region string, options ...interface{}) (aws.Endpoint, error) {
							if service == awsS3.ServiceID && region == config.S3Region {
								return aws.Endpoint{
									URL:           config.S3Endpoint,
									SigningRegion: config.S3Region,
								}, nil
							}
							return aws.Endpoint{}, &aws.EndpointNotFoundError{}
						},
					),
				}

				otelaws.AppendMiddlewares(&cfg.APIOptions)

				return awsS3.NewFromConfig(cfg), nil
			},
			func(config *ModuleConfig) *s3.ImageUrlConstructor {
				return s3.NewImageUrlConstructor(config.CdnHost, config.ProtectionSecret)
			},
			func(
				client *awsS3.Client,
				urlConstructor *s3.ImageUrlConstructor,
				logger *zap.Logger,
				config *ModuleConfig,
			) *s3.ImageStorage {
				return s3.NewImageStorage(client, urlConstructor, logger, config.S3BucketForImages)
			},

			func(urlConstructor *s3.ImageUrlConstructor) image.UrlConstructor {
				return urlConstructor
			},
			func(storage *s3.ImageStorage) image.Uploader {
				return storage
			},
		),
	)
}
