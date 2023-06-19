package dev

import (
	"musicd/internal/cli"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"dev",
		fx.Provide(
			NewStarter,

			cli.ProvideCommand(
				func(starter *Starter) *cobra.Command {
					root := &cobra.Command{Use: "dev"}

					root.AddCommand(starter.Command())

					return root
				},
			),
		),
	)
}
