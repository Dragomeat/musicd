package player

import (
	"musicd/internal/player/api"
	"musicd/internal/player/domain"
	"musicd/internal/player/storage"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"player",
		fx.Provide(
			api.NewPlayerResolver,
			api.NewPlayerTransformer,

			storage.NewPlayers,
			func(players *storage.Players) domain.Players {
				return players
			},
			func(pool *pgxpool.Pool) storage.DBTX {
				return pool
			},
			func(db storage.DBTX) *storage.Queries {
				return storage.New(db)
			},
		))
}
