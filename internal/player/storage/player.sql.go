// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: player.sql

package storage

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"musicd/internal/player/domain"
)

const createPlayer = `-- name: CreatePlayer :exec
insert into "player" (id, host_id, updated_at, created_at)
values ($1, $2, $3, $4)
`

type CreatePlayerParams struct {
	ID        uuid.UUID
	HostID    uuid.UUID
	UpdatedAt time.Time
	CreatedAt time.Time
}

func (q *Queries) CreatePlayer(ctx context.Context, arg CreatePlayerParams) error {
	_, err := q.db.Exec(ctx, createPlayer,
		arg.ID,
		arg.HostID,
		arg.UpdatedAt,
		arg.CreatedAt,
	)
	return err
}

const findPlayer = `-- name: FindPlayer :one
select id, host_id, current_track, queue, updated_at, created_at from "player" where id = $1
`

func (q *Queries) FindPlayer(ctx context.Context, id uuid.UUID) (Player, error) {
	row := q.db.QueryRow(ctx, findPlayer, id)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.HostID,
		&i.CurrentTrack,
		&i.Queue,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const findPlayerForUpdate = `-- name: FindPlayerForUpdate :one
select id, host_id, current_track, queue, updated_at, created_at from "player" where id = $1 for update
`

func (q *Queries) FindPlayerForUpdate(ctx context.Context, id uuid.UUID) (Player, error) {
	row := q.db.QueryRow(ctx, findPlayerForUpdate, id)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.HostID,
		&i.CurrentTrack,
		&i.Queue,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const updatePlayer = `-- name: UpdatePlayer :exec
update "player" set host_id = $1, current_track = $2, queue = $3, updated_at = $4 where id = $5
`

type UpdatePlayerParams struct {
	HostID       uuid.UUID
	CurrentTrack domain.NullCurrentTrack
	Queue        domain.Queue
	UpdatedAt    time.Time
	ID           uuid.UUID
}

func (q *Queries) UpdatePlayer(ctx context.Context, arg UpdatePlayerParams) error {
	_, err := q.db.Exec(ctx, updatePlayer,
		arg.HostID,
		arg.CurrentTrack,
		arg.Queue,
		arg.UpdatedAt,
		arg.ID,
	)
	return err
}
