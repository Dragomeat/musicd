package pgx

import (
	"context"
	"fmt"
	"musicd/internal/logger"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type ModuleConfig struct {
	PgDsn       string `mapstructure:"PGX_DSN"`
	SlowQueryMs int    `mapstructure:"PGX_SLOW_QUERY_LOGGING_LIMIT"`
}

func NewPgx(
	lc fx.Lifecycle,
	config *ModuleConfig,
	logger logger.Logger,
) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(config.PgDsn)
	if err != nil {
		return nil, fmt.Errorf("can`t parse pg config: %w", err)
	}

	poolConfig.ConnConfig.Tracer = otelpgx.NewTracer(
		otelpgx.WithIncludeQueryParameters(),
	)

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("can`t create pgx pool: %w", err)
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				ctx, cancel := context.WithTimeout(ctx, time.Second*5)
				defer cancel()

				err := pool.Ping(ctx)
				if err != nil {
					return fmt.Errorf("can`t connect to db: %w", err)
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				pool.Close()
				return nil
			},
		},
	)

	return pool, nil
}

func Module(config ModuleConfig) fx.Option {
	return fx.Module(
		"pgx",
		fx.Provide(
			NewPgx,
			func(viper *viper.Viper) (*ModuleConfig, error) {
				err := viper.Unmarshal(&config)
				if err != nil {
					return nil, err
				}
				return &config, nil
			},
		),
	)
}
