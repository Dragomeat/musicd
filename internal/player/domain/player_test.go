package domain_test

import (
	"math/rand"
	"musicd/internal/player/domain"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTrack(t *testing.T) {
	clock := clockwork.NewFakeClock()
	performerId := uuid.Must(uuid.NewV6())

	player, err := domain.NewPlayer(clock, performerId)
	require.NoError(t, err)

	assert.Equal(t, player.HostID, performerId)
	assert.Empty(t, player.Queue)
	assert.False(t, player.CurrentTrack.Valid)
	assert.Equal(t, player.CreatedAt, clock.Now())
	assert.Equal(t, player.UpdatedAt, clock.Now())
}

func TestAddTrackToQueue(t *testing.T) {
	clock := clockwork.NewFakeClock()
	performerId := uuid.Must(uuid.NewV6())
	trackId := uuid.Must(uuid.NewV6())
	player, err := domain.NewPlayer(clock, performerId)
	require.NoError(t, err)
	clock.Advance(5 * time.Second)

	err = player.AddTrackToQueue(clock, performerId, trackId)
	require.NoError(t, err)

	assert.False(t, player.CurrentTrack.Valid)
	assert.Contains(t, player.Queue, domain.QueuedTrack{TrackID: trackId, QueueedBy: performerId, QueuedAt: clock.Now()})
}

func TestRemoveTrackFromQueue(t *testing.T) {
	clock := clockwork.NewFakeClock()
	performerId := uuid.Must(uuid.NewV6())
	trackId := uuid.Must(uuid.NewV6())
	player, err := domain.NewPlayer(clock, performerId)
	require.NoError(t, err)
	clock.Advance(5 * time.Second)
	err = player.AddTrackToQueue(clock, performerId, trackId)
	require.NoError(t, err)
	clock.Advance(5 * time.Second)

	err = player.RemoveTrackFromQueue(performerId, trackId)
	require.NoError(t, err)

	assert.False(t, player.CurrentTrack.Valid)
	assert.Empty(t, player.Queue)
}

func TestMoveTrackInQueue(t *testing.T) {
	clock := clockwork.NewFakeClock()
	performerId := uuid.Must(uuid.NewV6())
	player, err := domain.NewPlayer(clock, performerId)
	require.NoError(t, err)
	clock.Advance(5 * time.Second)

	trackIds := make([]uuid.UUID, 10)
	for i := range trackIds {
		trackId := uuid.Must(uuid.NewV6())
		trackIds[i] = trackId

		err = player.AddTrackToQueue(clock, performerId, trackId)
		require.NoError(t, err)
	}

	currentPosition := rand.Intn(9)
	trackIdA := trackIds[currentPosition]
	var nextPosition int
	for n := true; n; n = currentPosition == nextPosition {
		nextPosition = rand.Intn(9)
	}
	trackIdB := trackIds[nextPosition]

	err = player.MoveTrackInQueue(performerId, trackIdA, nextPosition)
	require.NoError(t, err)

	assert.Equal(t, player.Queue[nextPosition].TrackID, trackIdA)
	assert.Equal(t, player.Queue[currentPosition].TrackID, trackIdB)
}

func TestStartPlayer(t *testing.T) {
	clock := clockwork.NewFakeClock()
	performerId := uuid.Must(uuid.NewV6())
	trackId := uuid.Must(uuid.NewV6())
	player, err := domain.NewPlayer(clock, performerId)
	require.NoError(t, err)
	clock.Advance(5 * time.Second)
	err = player.AddTrackToQueue(clock, performerId, trackId)
	require.NoError(t, err)
	clock.Advance(5 * time.Second)

	err = player.Start(clock, performerId)
	require.NoError(t, err)

	assert.Empty(t, player.Queue)
	assert.True(t, player.CurrentTrack.Valid)
	assert.Equal(t, player.CurrentTrack.CurrentTrack.TrackID, trackId)
	assert.True(t, player.CurrentTrack.CurrentTrack.Playing)
	assert.Equal(t, player.CurrentTrack.CurrentTrack.StartedAt, clock.Now())
	assert.False(t, player.CurrentTrack.CurrentTrack.PausedAt.Valid)
	assert.Equal(t, player.CurrentTrack.CurrentTrack.Cursor, time.Duration(0))
	assert.Equal(t, player.CurrentTrack.CurrentTrack.QueuedBy, performerId)
	assert.Equal(t, player.CurrentTrack.CurrentTrack.QueuedAt, clock.Now().Add(-5*time.Second))
}

func TestStopPlayer(t *testing.T) {
	clock := clockwork.NewFakeClock()
	performerId := uuid.Must(uuid.NewV6())
	trackId := uuid.Must(uuid.NewV6())
	player, err := domain.NewPlayer(clock, performerId)
	require.NoError(t, err)
	clock.Advance(5 * time.Second)
	err = player.AddTrackToQueue(clock, performerId, trackId)
	require.NoError(t, err)
	clock.Advance(5 * time.Second)
	err = player.Start(clock, performerId)
	require.NoError(t, err)
	clock.Advance(15 * time.Second)

	err = player.Stop(clock, performerId)
	require.NoError(t, err)

	assert.Empty(t, player.Queue)
	assert.True(t, player.CurrentTrack.Valid)
	assert.Equal(t, player.CurrentTrack.CurrentTrack.TrackID, trackId)
	assert.False(t, player.CurrentTrack.CurrentTrack.Playing)
	assert.Equal(t, player.CurrentTrack.CurrentTrack.StartedAt, clock.Now().Add(-15*time.Second))
	assert.True(t, player.CurrentTrack.CurrentTrack.PausedAt.Valid)
	assert.Equal(t, player.CurrentTrack.CurrentTrack.PausedAt.Time, clock.Now())
	assert.Equal(t, player.CurrentTrack.CurrentTrack.Cursor, time.Duration(0))
	assert.Equal(t, player.CurrentTrack.CurrentTrack.QueuedBy, performerId)
	assert.Equal(t, player.CurrentTrack.CurrentTrack.QueuedAt, clock.Now().Add(-20*time.Second))
}

func TestNextTrack(t *testing.T) {
	clock := clockwork.NewFakeClock()
	performerId := uuid.Must(uuid.NewV6())
	player, err := domain.NewPlayer(clock, performerId)
	require.NoError(t, err)
	clock.Advance(5 * time.Second)
	trackIds := make([]uuid.UUID, 10)
	for i := range trackIds {
		trackId := uuid.Must(uuid.NewV6())
		trackIds[i] = trackId

		err = player.AddTrackToQueue(clock, performerId, trackId)
		require.NoError(t, err)
	}
	err = player.Start(clock, performerId)
	require.NoError(t, err)
	require.Equal(t, player.Queue[0].TrackID, trackIds[1])
	clock.Advance(10 * time.Second)

	err = player.NextTrack(clock, performerId)
	require.NoError(t, err)

	assert.Equal(t, player.Queue[0].TrackID, trackIds[2])
	assert.True(t, player.CurrentTrack.Valid)
	assert.Equal(t, player.CurrentTrack.CurrentTrack.TrackID, trackIds[1])
	assert.True(t, player.CurrentTrack.CurrentTrack.Playing)
	assert.Equal(t, player.CurrentTrack.CurrentTrack.StartedAt, clock.Now())
	assert.False(t, player.CurrentTrack.CurrentTrack.PausedAt.Valid)
	assert.Equal(t, player.CurrentTrack.CurrentTrack.Cursor, time.Duration(0))
	assert.Equal(t, player.CurrentTrack.CurrentTrack.QueuedBy, performerId)
	assert.Equal(t, player.CurrentTrack.CurrentTrack.QueuedAt, clock.Now().Add(-10*time.Second))
}

func TestPreviousTrack(t *testing.T) {
	clock := clockwork.NewFakeClock()
	performerId := uuid.Must(uuid.NewV6())
	player, err := domain.NewPlayer(clock, performerId)
	require.NoError(t, err)
	clock.Advance(5 * time.Second)
	trackIds := make([]uuid.UUID, 10)
	for i := range trackIds {
		trackId := uuid.Must(uuid.NewV6())
		trackIds[i] = trackId

		err = player.AddTrackToQueue(clock, performerId, trackId)
		require.NoError(t, err)
	}
	err = player.Start(clock, performerId)
	require.NoError(t, err)
	require.Equal(t, player.Queue[0].TrackID, trackIds[1])
	clock.Advance(10 * time.Second)

	err = player.PreviousTrack(clock, performerId)
	require.NoError(t, err)

	assert.Equal(t, player.Queue[0].TrackID, trackIds[0])
	assert.True(t, player.CurrentTrack.Valid)
	assert.Equal(t, player.CurrentTrack.CurrentTrack.TrackID, trackIds[9])
	assert.True(t, player.CurrentTrack.CurrentTrack.Playing)
	assert.Equal(t, player.CurrentTrack.CurrentTrack.StartedAt, clock.Now())
	assert.False(t, player.CurrentTrack.CurrentTrack.PausedAt.Valid)
	assert.Equal(t, player.CurrentTrack.CurrentTrack.Cursor, time.Duration(0))
	assert.Equal(t, player.CurrentTrack.CurrentTrack.QueuedBy, performerId)
	assert.Equal(t, player.CurrentTrack.CurrentTrack.QueuedAt, clock.Now().Add(-10*time.Second))
}

func TestSeekTrack(t *testing.T) {
	clock := clockwork.NewFakeClock()
	performerId := uuid.Must(uuid.NewV6())
	player, err := domain.NewPlayer(clock, performerId)
	require.NoError(t, err)
	clock.Advance(5 * time.Second)
	trackId := uuid.Must(uuid.NewV6())
	err = player.AddTrackToQueue(clock, performerId, trackId)
	require.NoError(t, err)
	err = player.Start(clock, performerId)
	require.NoError(t, err)
	clock.Advance(10 * time.Second)

	err = player.Seek(performerId, time.Duration(5*time.Second))
	require.NoError(t, err)

	assert.True(t, player.CurrentTrack.Valid)
	assert.Equal(t, player.CurrentTrack.CurrentTrack.TrackID, trackId)
	assert.True(t, player.CurrentTrack.CurrentTrack.Playing)
	assert.Equal(t, player.CurrentTrack.CurrentTrack.StartedAt, clock.Now().Add(-10*time.Second))
	assert.False(t, player.CurrentTrack.CurrentTrack.PausedAt.Valid)
	assert.Equal(t, player.CurrentTrack.CurrentTrack.Cursor, time.Duration(5*time.Second))
	assert.Equal(t, player.CurrentTrack.CurrentTrack.QueuedBy, performerId)
	assert.Equal(t, player.CurrentTrack.CurrentTrack.QueuedAt, clock.Now().Add(-10*time.Second))
}
