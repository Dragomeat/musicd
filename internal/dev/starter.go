package dev

import (
	"musicd/internal/chi"
	"musicd/internal/temporal"

	"github.com/spf13/cobra"
)

type Starter struct {
	serve  *chi.Serve
	worker *temporal.Worker
}

func NewStarter(serve *chi.Serve, worker *temporal.Worker) *Starter {
	return &Starter{
		serve:  serve,
		worker: worker,
	}
}

func (w *Starter) Command() *cobra.Command {
	return &cobra.Command{
		Use:  "start",
		Args: cobra.NoArgs,
		RunE: w.Run,
	}
}

func (w *Starter) Run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	go func() {
		downloaderCmd := w.worker.Command()
		downloaderCmd.SetArgs([]string{})
		downloaderCmd.Flags().Set(temporal.FlagQueue, "track")
		downloaderCmd.ExecuteContext(ctx)
	}()

	serveCmd := w.serve.Command()
	serveCmd.SetArgs([]string{})
	return serveCmd.ExecuteContext(ctx)
}
