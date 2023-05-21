package domain

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"gopkg.in/guregu/null.v4"
)

type Artist struct {
	ID    uuid.UUID
	Title string
}

type Album struct {
	ID         uuid.UUID
	Title      string
	ArtistIDs  []uuid.UUID
	TrackIDs   []uuid.UUID
	Cover      string
	ReleasedAt time.Time
}

type Playlist struct {
	ID       uuid.UUID
	Title    string
	TrackIDs []uuid.UUID
}

type Files []Audio

func (f Files) Value() (driver.Value, error) {
	return json.Marshal(f)
}

type ExternalSource struct {
	value string
}

func (e ExternalSource) String() string {
	return e.value
}

func (e *ExternalSource) UnmarshalJSON(data []byte) error {
	var err error
	*e, err = NewExternalSourceFromString(string(data))
	return err
}

func (e ExternalSource) MarshalJSON() ([]byte, error) {
	return []byte(`"` + e.value + `"`), nil
}

var (
	ExternalSourceUnknown = ExternalSource{""}
	ExternalSourceYoutube = ExternalSource{"youtube"}
)

func NewExternalSourceFromString(value string) (ExternalSource, error) {
	switch value {
	case ExternalSourceYoutube.value:
		return ExternalSourceYoutube, nil
	default:
		return ExternalSourceUnknown, nil
	}
}

type MIMEType struct {
	value string
}

func (m MIMEType) String() string {
	return m.value
}

func (m *MIMEType) UnmarshalJSON(data []byte) error {
	var err error
	*m, err = NewMIMETypeFromString(string(data))
	return err
}

func (m MIMEType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + m.value + `"`), nil
}

var (
	MIMETypeAudioUnknown = MIMEType{""}
	MIMETypeAudioMp3     = MIMEType{"audio/mp3"}
	MIMETypeAudioMp4     = MIMEType{"audio/mp4"}
)

func NewMIMETypeFromString(value string) (MIMEType, error) {
	switch value {
	case MIMETypeAudioMp3.value:
		return MIMETypeAudioMp3, nil
	case MIMETypeAudioMp4.value:
		return MIMETypeAudioMp4, nil
	default:
		return MIMETypeAudioUnknown, nil
	}
}

type Audio struct {
	ExternalId     string         `json:"externalId"`
	ExternalSource ExternalSource `json:"externalSource"`
	Sha256         string         `json:"sha256"`
	Type           MIMEType       `json:"type"`
	Bitrate        int            `json:"bitrate"`
}

type Track struct {
	ID         uuid.UUID
	ArtistID   uuid.NullUUID
	AlbumID    uuid.NullUUID
	Title      string
	Duration   int
	Files      Files
	UploadedBy uuid.UUID
	UploadedAt time.Time
}

func NewTrack(
	name string,
	duration int,
	files []Audio,
	uploadedBy uuid.UUID,
) (Track, error) {
	id, err := uuid.NewV6()
	if err != nil {
		return Track{}, err
	}

	return Track{
		ID:         id,
		Title:      name,
		Duration:   duration,
		Files:      files,
		UploadedBy: uploadedBy,
		UploadedAt: time.Now(),
	}, nil
}

type TrackEdge struct {
	Cursor string
	Node   *Track
}

type Tracks interface {
	FindTracks(ctx context.Context, trackIds []uuid.UUID) (map[uuid.UUID]*Track, error)
	PaginateTracks(ctx context.Context, first int, after null.String, before null.String) ([]TrackEdge, error)
	Create(ctx context.Context, track Track) error
}

type NullTracks struct{}

func (NullTracks) FindTracks(ctx context.Context, trackIds []uuid.UUID) (map[uuid.UUID]*Track, error) {
	return nil, nil
}

func (NullTracks) PaginateTracks(ctx context.Context, first int, after null.String, before null.String) ([]TrackEdge, error) {
	return nil, nil
}

func (NullTracks) Create(ctx context.Context, track Track) error {
	return nil
}
