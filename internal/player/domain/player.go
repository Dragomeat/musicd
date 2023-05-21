package domain

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"musicd/internal/errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jonboulle/clockwork"
	"gopkg.in/guregu/null.v4"
)

type PlayerNotFoundError struct {
	ID uuid.UUID
}

func (e PlayerNotFoundError) Code() errors.ErrorCode {
	return "player.notFound"
}

func (e PlayerNotFoundError) Message() string {
	return "Player not found"
}

func (e PlayerNotFoundError) Error() string {
	return fmt.Sprintf("player not found: %s", e.ID)
}

type Players interface {
	// Get returns a player by ID. If the player does not exist, it returns a [PlayerNotFoundError].
	Get(ctx context.Context, id uuid.UUID) (*Player, error)
	Update(ctx context.Context, id uuid.UUID, updateFunc func(ctx context.Context, player *Player) (*Player, error)) (*Player, error)
	Create(ctx context.Context, player *Player) error
}

type NullCurrentTrack struct {
	CurrentTrack CurrentTrack
	Valid        bool
}

func (t *NullCurrentTrack) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	t.Valid = true
	return json.Unmarshal([]byte(value.(string)), &t.CurrentTrack)
}

func (t NullCurrentTrack) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}

	return json.Marshal(t.CurrentTrack)
}

type CurrentTrack struct {
	TrackID   uuid.UUID     `json:"trackId"`
	Playing   bool          `json:"playing"`
	Cursor    time.Duration `json:"cursor"`
	StartedAt time.Time     `json:"startedAt"`
	PausedAt  null.Time     `json:"pausedAt"`
	QueuedBy  uuid.UUID     `json:"queuedBy"`
	QueuedAt  time.Time     `json:"queuedAt"`
}

func (t *CurrentTrack) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), t)
}

func (t CurrentTrack) Value() (driver.Value, error) {
	return json.Marshal(t)
}

type QueuedTrack struct {
	TrackID   uuid.UUID `json:"trackId"`
	QueueedBy uuid.UUID `json:"queuedBy"`
	QueuedAt  time.Time `json:"queuedAt"`
}

type Queue []QueuedTrack

func (q *Queue) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal([]byte(value.(string)), q)
}

func (q *Queue) Value() (driver.Value, error) {
	return json.Marshal(q)
}

type Player struct {
	ID           uuid.UUID
	HostID       uuid.UUID
	CurrentTrack NullCurrentTrack
	Queue        Queue
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

func NewPlayer(clock clockwork.Clock, hostID uuid.UUID) (*Player, error) {
	id, err := uuid.NewV6()
	if err != nil {
		return nil, err
	}

	now := clock.Now()

	return &Player{
		ID:        id,
		HostID:    hostID,
		Queue:     make(Queue, 0),
		UpdatedAt: now,
		CreatedAt: now,
	}, nil
}

func (p *Player) AddTrackToQueue(
	clock clockwork.Clock,
	performerID uuid.UUID,
	trackID uuid.UUID,
) error {
	for _, queuedTrack := range p.Queue {
		if queuedTrack.TrackID == trackID {
			return fmt.Errorf("track already in queue") // TODO: create user friendly error
		}
	}
	p.Queue = append(p.Queue, QueuedTrack{TrackID: trackID, QueueedBy: performerID, QueuedAt: clock.Now()})

	return nil
}

func (p *Player) RemoveTrackFromQueue(performerID uuid.UUID, trackID uuid.UUID) error {
	for i, queuedTrack := range p.Queue {
		if queuedTrack.TrackID == trackID {
			p.Queue = append(p.Queue[:i], p.Queue[i+1:]...)
			break
		}
	}
	return nil
}

func (p *Player) MoveTrackInQueue(performerID uuid.UUID, trackID uuid.UUID, position int) error {
	currentPosition := -1
	for i, queuedTrack := range p.Queue {
		if queuedTrack.TrackID == trackID {
			currentPosition = i
			break
		}
	}
	if currentPosition == -1 {
		return fmt.Errorf("track not found")
	}
	if currentPosition == position {
		return nil
	}

	p.Queue[currentPosition], p.Queue[position] = p.Queue[position], p.Queue[currentPosition]

	return nil
}

func (p *Player) Start(clock clockwork.Clock, performerID uuid.UUID) error {
	if p.CurrentTrack.CurrentTrack.Playing {
		return nil
	}

	if !p.CurrentTrack.Valid {
		if len(p.Queue) == 0 {
			return fmt.Errorf("no track to play") // TODO: create user friendly error
		}

		var nextTrack QueuedTrack
		nextTrack, p.Queue = p.Queue[0], p.Queue[1:]
		p.CurrentTrack = NullCurrentTrack{
			CurrentTrack{
				TrackID:   nextTrack.TrackID,
				StartedAt: clock.Now(),
				QueuedBy:  nextTrack.QueueedBy,
				QueuedAt:  nextTrack.QueuedAt,
			},
			true,
		}
	}

	p.CurrentTrack.CurrentTrack.Playing = true
	p.CurrentTrack.CurrentTrack.PausedAt = null.Time{}

	return nil
}

func (p *Player) Stop(clock clockwork.Clock, performerID uuid.UUID) error {
	if !p.CurrentTrack.CurrentTrack.Playing {
		return nil
	}

	p.CurrentTrack.CurrentTrack.Playing = false
	p.CurrentTrack.CurrentTrack.PausedAt = null.TimeFrom(clock.Now())

	return nil
}

func (p *Player) PreviousTrack(clock clockwork.Clock, performerID uuid.UUID) error {
	if len(p.Queue) == 0 {
		return fmt.Errorf("no previous track") // TODO: create user friendly error
	}

	currentTrack := p.CurrentTrack
	var nextTrack QueuedTrack
	nextTrack, p.Queue = p.Queue[len(p.Queue)-1], p.Queue[:len(p.Queue)-1]
	p.CurrentTrack = NullCurrentTrack{
		CurrentTrack{
			TrackID:   nextTrack.TrackID,
			Playing:   currentTrack.CurrentTrack.Playing,
			StartedAt: clock.Now(),
			QueuedBy:  nextTrack.QueueedBy,
			QueuedAt:  nextTrack.QueuedAt,
		},
		true,
	}

	if currentTrack.Valid {
		p.Queue = append(
			[]QueuedTrack{
				{
					TrackID:   currentTrack.CurrentTrack.TrackID,
					QueueedBy: currentTrack.CurrentTrack.QueuedBy,
					QueuedAt:  currentTrack.CurrentTrack.QueuedAt,
				},
			},
			p.Queue...,
		)
	}
	return nil
}

func (p *Player) NextTrack(clock clockwork.Clock, performerID uuid.UUID) error {
	if !p.CurrentTrack.Valid {
		return fmt.Errorf("no track to skip") // TODO: create user friendly error
	}
	if len(p.Queue) == 0 {
		return fmt.Errorf("no previous track") // TODO: create user friendly error
	}

	currentTrack := p.CurrentTrack.CurrentTrack
	var nextTrack QueuedTrack
	nextTrack, p.Queue = p.Queue[0], p.Queue[1:]
	p.CurrentTrack = NullCurrentTrack{
		CurrentTrack{
			TrackID:   nextTrack.TrackID,
			Playing:   currentTrack.Playing,
			StartedAt: clock.Now(),
			QueuedBy:  nextTrack.QueueedBy,
			QueuedAt:  nextTrack.QueuedAt,
		},
		true,
	}
	p.Queue = append(
		p.Queue,
		QueuedTrack{
			TrackID:   currentTrack.TrackID,
			QueueedBy: currentTrack.QueuedBy,
			QueuedAt:  currentTrack.QueuedAt,
		},
	)

	return nil
}

func (p *Player) Seek(performerID uuid.UUID, cursor time.Duration) error {
	if !p.CurrentTrack.Valid {
		return fmt.Errorf("no track to seek") // TODO: create user friendly error
	}

	p.CurrentTrack.CurrentTrack.Cursor = cursor

	return nil
}

func (p *Player) ClearQueue(performerID uuid.UUID) error {
	p.Queue = make(Queue, 0)
	return nil
}
