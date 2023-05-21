package api

import (
	"context"
	"musicd/graph/model"
	"musicd/internal/auth"
	"musicd/internal/errors"
	"musicd/internal/player/domain"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jonboulle/clockwork"
)

type PlayerResolver struct {
	clock             clockwork.Clock
	players           domain.Players
	playerTransformer *PlayerTransformer
}

func NewPlayerResolver(
	clock clockwork.Clock,
	players domain.Players,
	playerTransformer *PlayerTransformer,
) *PlayerResolver {
	return &PlayerResolver{
		clock:             clock,
		players:           players,
		playerTransformer: playerTransformer,
	}
}

func (r *PlayerResolver) GetPlayer(ctx context.Context, id uuid.UUID) (*model.Player, error) {
	player, err := r.players.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.playerTransformer.Transform(player)
}

func (r *PlayerResolver) CreatePlayer(ctx context.Context) (*model.Player, error) {
	performer := auth.PerformerFromContext(ctx)
	if !performer.Valid {
		return nil, errors.WithStack(auth.ErrUnathenticated)
	}

	player, err := domain.NewPlayer(r.clock, performer.Value.Id)
	if err != nil {
		return nil, err
	}

	err = r.players.Create(ctx, player)
	if err != nil {
		return nil, err
	}

	return r.playerTransformer.Transform(player)
}

func (r *PlayerResolver) QueueTrack(ctx context.Context, playerId uuid.UUID, trackId uuid.UUID) (*model.Player, error) {
	performer := auth.PerformerFromContext(ctx)
	if !performer.Valid {
		return nil, errors.WithStack(auth.ErrUnathenticated)
	}
	player, err := r.players.Update(
		ctx,
		playerId,
		func(ctx context.Context, player *domain.Player) (*domain.Player, error) {
			err := player.AddTrackToQueue(r.clock, performer.Value.Id, trackId)
			if err != nil {
				return nil, err
			}

			return player, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return r.playerTransformer.Transform(player)
}

func (r *PlayerResolver) RemoveTrackFromQueue(ctx context.Context, playerId uuid.UUID, trackId uuid.UUID) (*model.Player, error) {
	performer := auth.PerformerFromContext(ctx)
	if !performer.Valid {
		return nil, errors.WithStack(auth.ErrUnathenticated)
	}
	player, err := r.players.Update(
		ctx,
		playerId,
		func(ctx context.Context, player *domain.Player) (*domain.Player, error) {
			err := player.RemoveTrackFromQueue(performer.Value.Id, trackId)
			if err != nil {
				return nil, err
			}

			return player, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return r.playerTransformer.Transform(player)
}

func (r *PlayerResolver) Start(ctx context.Context, playerId uuid.UUID) (*model.Player, error) {
	performer := auth.PerformerFromContext(ctx)
	if !performer.Valid {
		return nil, errors.WithStack(auth.ErrUnathenticated)
	}
	player, err := r.players.Update(
		ctx,
		playerId,
		func(ctx context.Context, player *domain.Player) (*domain.Player, error) {
			err := player.Start(r.clock, performer.Value.Id)
			if err != nil {
				return nil, err
			}

			return player, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return r.playerTransformer.Transform(player)
}

func (r *PlayerResolver) Stop(ctx context.Context, playerId uuid.UUID) (*model.Player, error) {
	performer := auth.PerformerFromContext(ctx)
	if !performer.Valid {
		return nil, errors.WithStack(auth.ErrUnathenticated)
	}
	player, err := r.players.Update(
		ctx,
		playerId,
		func(ctx context.Context, player *domain.Player) (*domain.Player, error) {
			err := player.Stop(r.clock, performer.Value.Id)
			if err != nil {
				return nil, err
			}

			return player, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return r.playerTransformer.Transform(player)
}

func (r *PlayerResolver) SeekTo(ctx context.Context, playerId uuid.UUID, positionInSeconds int) (*model.Player, error) {
	performer := auth.PerformerFromContext(ctx)
	if !performer.Valid {
		return nil, errors.WithStack(auth.ErrUnathenticated)
	}
	player, err := r.players.Update(
		ctx,
		playerId,
		func(ctx context.Context, player *domain.Player) (*domain.Player, error) {
			err := player.Seek(performer.Value.Id, time.Duration(positionInSeconds)*time.Second)
			if err != nil {
				return nil, err
			}

			return player, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return r.playerTransformer.Transform(player)
}

func (r *PlayerResolver) NextTrack(ctx context.Context, playerId uuid.UUID) (*model.Player, error) {
	performer := auth.PerformerFromContext(ctx)
	if !performer.Valid {
		return nil, errors.WithStack(auth.ErrUnathenticated)
	}
	player, err := r.players.Update(
		ctx,
		playerId,
		func(ctx context.Context, player *domain.Player) (*domain.Player, error) {
			err := player.NextTrack(r.clock, performer.Value.Id)
			if err != nil {
				return nil, err
			}

			return player, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return r.playerTransformer.Transform(player)
}

func (r *PlayerResolver) PreviousTrack(ctx context.Context, playerId uuid.UUID) (*model.Player, error) {
	performer := auth.PerformerFromContext(ctx)
	if !performer.Valid {
		return nil, errors.WithStack(auth.ErrUnathenticated)
	}
	player, err := r.players.Update(
		ctx,
		playerId,
		func(ctx context.Context, player *domain.Player) (*domain.Player, error) {
			err := player.PreviousTrack(r.clock, performer.Value.Id)
			if err != nil {
				return nil, err
			}

			return player, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return r.playerTransformer.Transform(player)
}

type PlayerTransformer struct{}

func NewPlayerTransformer() *PlayerTransformer {
	return &PlayerTransformer{}
}

func (t *PlayerTransformer) Transform(player *domain.Player) (*model.Player, error) {
	tracksInQueue := 0
	queue := make([]*model.QueuedTrack, len(player.Queue))
	for i, queuedTrack := range player.Queue {
		queue[i] = &model.QueuedTrack{
			Track: &model.Track{
				ID: queuedTrack.TrackID,
			},
			QueuedBy: &model.Person{ID: queuedTrack.QueueedBy},
			QueuedAt: queuedTrack.QueuedAt,
		}
	}
	var currentTrack *model.CurrentTrack
	if player.CurrentTrack.Valid {
		ct := player.CurrentTrack.CurrentTrack
		currentTrack = &model.CurrentTrack{
			Track: &model.Track{
				ID: ct.TrackID,
			},
			QueuedBy:          &model.Person{ID: ct.QueuedBy},
			Playing:           ct.Playing,
			PositionInSeconds: int(ct.Cursor.Seconds()),
		}
	}
	return &model.Player{
		ID:            player.ID,
		Host:          &model.Person{ID: player.HostID},
		CurrentTrack:  currentTrack,
		Queue:         queue,
		TracksInQueue: tracksInQueue,
	}, nil
}
