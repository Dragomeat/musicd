package storage

import (
	"context"
	"musicd/internal/errors"
	"musicd/internal/player/domain"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Players struct {
	db      *pgxpool.Pool
	queries *Queries
}

func NewPlayers(db *pgxpool.Pool, queries *Queries) *Players {
	return &Players{db: db, queries: queries}
}

func (p *Players) Get(ctx context.Context, id uuid.UUID) (*domain.Player, error) {
	player, err := p.queries.FindPlayer(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.WithStack(domain.PlayerNotFoundError{ID: id})
		}
		return nil, err
	}

	return p.mapDtoToDomain(player), nil
}

func (p *Players) Update(ctx context.Context, id uuid.UUID, updateFunc func(ctx context.Context, player *domain.Player) (*domain.Player, error)) (*domain.Player, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := p.queries.WithTx(tx)
	rawPlayer, err := q.FindPlayerForUpdate(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.WithStack(domain.PlayerNotFoundError{ID: id})
		}

		return nil, err
	}

	player := p.mapDtoToDomain(rawPlayer)

	player, err = updateFunc(ctx, player)
	if err != nil {
		return nil, err
	}

	err = q.UpdatePlayer(
		ctx,
		UpdatePlayerParams{
			ID:           player.ID,
			HostID:       player.HostID,
			CurrentTrack: player.CurrentTrack,
			Queue:        player.Queue,
			UpdatedAt:    player.UpdatedAt,
		},
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func (p *Players) Create(ctx context.Context, player *domain.Player) error {
	return p.queries.CreatePlayer(ctx, CreatePlayerParams{
		ID:        player.ID,
		HostID:    player.HostID,
		UpdatedAt: player.UpdatedAt,
		CreatedAt: player.CreatedAt,
	})
}

func (p *Players) mapDtoToDomain(player Player) *domain.Player {
	return &domain.Player{
		ID:           player.ID,
		HostID:       player.HostID,
		CurrentTrack: player.CurrentTrack,
		Queue:        player.Queue,
		UpdatedAt:    player.UpdatedAt,
		CreatedAt:    player.CreatedAt,
	}
}
