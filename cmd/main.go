package main

import (
	"musicd/graph"
	"musicd/internal/auth"
	"musicd/internal/chi"
	"musicd/internal/cli"
	"musicd/internal/config"
	"musicd/internal/dev"
	"musicd/internal/errors"
	"musicd/internal/graphql"
	"musicd/internal/http"
	"musicd/internal/http/middleware"
	"musicd/internal/logger"
	"musicd/internal/media"
	"musicd/internal/pgx"
	"musicd/internal/player"
	"musicd/internal/temporal"
	"musicd/internal/tracing"
	"musicd/internal/track"
	oHttp "net/http"
	"reflect"
	"time"

	"github.com/ggicci/httpin"
	oChi "github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofrs/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	chiCfg := chi.ModuleParams{
		Configure: func(
			router oChi.Router,
			errorHandler *http.ErrorHandler,
			requestIdMiddleware *middleware.RequestIdMiddleware,
			correlationIdMiddleware *middleware.CorrelationIdMiddleware,
			authMiddleware *auth.Middleware,
		) {
			m := http.Chain(requestIdMiddleware, correlationIdMiddleware)
			router.Use(func(handler oHttp.Handler) oHttp.Handler {
				return errorHandler.Wrap(m.Next(http.FromHttpHandler(handler)))
			})
			router.Use(chiMiddleware.SetHeader("Content-Type", "application/json"))
			router.Use(chiMiddleware.CleanPath)
			router.Use(chiMiddleware.RealIP)
			router.Use(chiMiddleware.Timeout(5 * time.Second))
			router.Use(chiMiddleware.NewCompressor(4, "application/json").Handler)
			router.Use(cors.Handler(cors.Options{AllowedHeaders: []string{"*"}, AllowCredentials: true}))
			router.Use(func(handler oHttp.Handler) oHttp.Handler {
				return errorHandler.Wrap(authMiddleware.Next(http.FromHttpHandler(handler)))
			})
			router.NotFound(func(w oHttp.ResponseWriter, req *oHttp.Request) {
				errorHandler.Handle(w, req, errors.NewUserError("http.notFound", "Not Found"))
			})
		},
	}

	app := fx.New(
		config.NewModule(),
		errors.Module(),
		logger.NewModule(),
		tracing.NewModule(),
		cli.Module(),
		chi.Module(chiCfg),
		pgx.Module(pgx.ModuleConfig{}),
		graph.Module(),
		graphql.Module(),
		temporal.Module(),
		dev.NewModule(),
		auth.Module(),
		media.NewModule(),
		player.NewModule(),
		track.NewModule(),
		fx.Provide(
			middleware.NewRequestIdMiddleware,
			middleware.NewCorrelationIdMiddleware,
			cli.ProvideRoot(
				func(shutdowner fx.Shutdowner) *cobra.Command {
					root := &cobra.Command{
						Use:   "musicd",
						Short: "Musicd is music server",
						Long:  `Musicd is place to store and listen your music`,
						Args:  cobra.NoArgs,
					}

					root.InitDefaultHelpCmd()
					root.InitDefaultHelpFlag()
					root.InitDefaultVersionFlag()

					return root
				},
			),
			func() *resource.Resource {
				r, _ := resource.Merge(
					resource.Default(),
					resource.NewWithAttributes(
						semconv.SchemaURL,
						semconv.ServiceName("musicd"),
						semconv.ServiceVersion("v0.1.0"),
						attribute.String("environment", "local"),
					),
				)
				return r
			},
			func() clockwork.Clock {
				return clockwork.NewRealClock()
			},
		),
		fx.Invoke(cli.Start),
		fx.WithLogger(
			func(logger *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: logger}
			},
		),
	)
	app.Run()
}

func init() {
	httpin.UseGochiURLParam("path", oChi.URLParam)
	httpin.RegisterTypeDecoder(
		reflect.TypeOf(uuid.UUID{}),
		httpin.ValueTypeDecoderFunc(
			func(s string) (interface{}, error) {
				if s == "" {
					return "", nil
				}

				err := validation.Validate(s, is.UUID)
				if err != nil {
					return nil, err
				}

				return uuid.FromString(s)
			},
		),
	)
	httpin.RegisterDirectiveExecutor(
		"ctx",
		httpin.DirectiveExecutorFunc(func(ctx *httpin.DirectiveContext) error {
			ctx.Value.Elem().Set(reflect.ValueOf(ctx.Context))

			return nil
		}),
		nil,
	)
}
