package chi

import (
	"context"
	"fmt"
	"musicd/internal/cli"
	"musicd/internal/http"
	"musicd/internal/logger"
	netHttp "net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/fx"
)

type Serve struct {
	lc           fx.Lifecycle
	runner       *cli.Runner
	config       *Config
	router       chi.Router
	routes       *http.Routes
	errorHandler *http.ErrorHandler
	errChannel   chan<- error
	logger       logger.Logger
}

type ServeParams struct {
	fx.In

	Lc           fx.Lifecycle
	Runner       *cli.Runner
	Config       *Config
	Router       chi.Router
	Routes       *http.Routes
	ErrorHandler *http.ErrorHandler
	ErrChannel   chan<- error `name:"errors.channel"`
	Logger       logger.Logger
}

func NewServe(params ServeParams) *Serve {
	return &Serve{
		lc:           params.Lc,
		runner:       params.Runner,
		config:       params.Config,
		router:       params.Router,
		routes:       params.Routes,
		errorHandler: params.ErrorHandler,
		errChannel:   params.ErrChannel,
		logger:       params.Logger,
	}
}

func (s *Serve) Command() *cobra.Command {
	return &cobra.Command{
		Use:  "serve",
		Args: cobra.NoArgs,
		Run:  s.Run,
	}
}

func (s *Serve) Run(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()

	server := &netHttp.Server{
		Addr:    s.config.Address,
		Handler: s.router,
		// TODO: Add logger
		// ErrorLog: logger,
	}

	for _, route := range s.routes.List() {
		s.logger.Info(ctx, "registering route", logger.Field("method", route.Method), logger.Field("path", route.Path))

		handler := route.Handler

		s.router.Method(
			route.Method,
			route.Path,
			otelhttp.NewHandler(s.errorHandler.Wrap(handler), fmt.Sprintf("%s %s", route.Method, route.Path)),
		)
	}

	s.runner.Run(ctx, func(ctx context.Context) {
		errChannel := make(chan error)
		go func() {
			s.logger.Info(ctx, "http server is starting")

			err := server.ListenAndServe()
			if err != nil {
				errChannel <- err
			}
		}()

		timer := time.NewTimer(1 * time.Second)

		for {
			select {
			case <-ctx.Done():
				s.logger.Info(ctx, "http server is stopping")

				err := server.Shutdown(ctx)
				if err != nil {
					s.errChannel <- err
					return
				}

				s.logger.Info(ctx, "http server has stopped")
				return
			case err := <-errChannel:
				if err == netHttp.ErrServerClosed {
					return
				}

				s.errChannel <- fmt.Errorf("http server has failed to run: %w", err)
				return
			case <-timer.C:
				s.logger.Info(
					ctx,
					"http server has started",
					logger.Field("address", s.config.Address),
				)
			}
		}
	})
}
