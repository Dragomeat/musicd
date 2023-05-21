package temporal

import (
	"fmt"
	"musicd/internal/cli"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentelemetry"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/worker"
	"go.uber.org/fx"
)

type Config struct {
	Address string `mapstructure:"TEMPORAL_ADDRESS"`
}

func NewConfig(viper *viper.Viper) (*Config, error) {
	config := &Config{
		Address: "localhost:7233",
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode temporal config: %w", err)
	}

	return config, nil
}

type Registerer interface {
	Register(worker.Registry)
}

func Provide[T Registerer](register interface{}) fx.Option {
	return fx.Provide(
		register,
		fx.Annotate(
			func(a T) T { return a },
			fx.As(new(Registerer)),
			fx.ResultTags(`group:"temporal.registerers"`),
		),
	)
}

func Module() fx.Option {
	return fx.Module(
		"temporal",
		fx.Provide(
			NewConfig,
			NewLogger,

			func(
				config *Config,
				logger *Logger,
			) (client.Client, error) {
				tracingInterceptor, err := opentelemetry.NewTracingInterceptor(opentelemetry.TracerOptions{})
				if err != nil {
					return nil, fmt.Errorf("unable to create tracing interceptor: %w", err)
				}

				opts := client.Options{
					HostPort:     config.Address,
					Logger:       logger,
					Interceptors: []interceptor.ClientInterceptor{tracingInterceptor},
				}

				return client.NewLazyClient(opts)
			},

			NewWorker,

			cli.ProvideCommand(
				func(worker *Worker) *cobra.Command {
					root := &cobra.Command{Use: "temporal"}

					root.AddCommand(worker.Command())

					return root
				},
			),
		),
	)
}
