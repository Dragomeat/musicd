package graphql

import (
	"context"
	"fmt"
	"musicd/internal/errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/ravilushqa/otelgqlgen"
	"github.com/spf13/viper"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type PlaygroundConfig struct {
	Enabled bool `mapstructure:"GRAPHQL_PLAYGROUND_ENABLED"`
	Path    string
}

type Config struct {
	Path       string
	Playground PlaygroundConfig `mapstructure:",squash"`
}

func NewConfig(viper *viper.Viper) (*Config, error) {
	config := &Config{
		Path: "/graphql",
		Playground: PlaygroundConfig{
			Enabled: true,
			Path:    "/playground",
		},
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode graphql config: %w", err)
	}

	return config, nil
}

type UserError interface {
	ToUserError() map[string]interface{}
}

func NewGraphqlServer(
	schema graphql.ExecutableSchema,
	loadersInitializer *LoadersInitializer,
	errorHandler *errors.ErrorHandler,
) *handler.Server {
	srv := handler.New(schema)

	srv.Use(loadersInitializer)
	srv.Use(otelgqlgen.Middleware())

	srv.SetRecoverFunc(func(ctx context.Context, p any) error {
		return fmt.Errorf("panic: %v", p)
	})

	srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		var gqlErr *gqlerror.Error
		path := graphql.GetPath(ctx)
		if errors.As(err, &gqlErr) {
			if gqlErr.Path == nil {
				gqlErr.Path = path
			} else {
				path = gqlErr.Path
			}

			originalErr := gqlErr.Unwrap()
			if originalErr == nil {
				return gqlErr
			}

			err = originalErr
		}

		code := errors.Code(err)
		message := errors.Message(err)
		extra := errors.Extra(err)
		if extra == nil {
			extra = make(map[string]any)
		}
		extra["code"] = code
		// if !uErr.DontHandle {
		errorHandler.Handle(ctx, err)
		// }

		return &gqlerror.Error{
			Message:    message,
			Path:       path,
			Extensions: extra,
		}
	})

	return srv
}
