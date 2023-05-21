// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type CurrentTrack struct {
	Track             *Track  `json:"track"`
	QueuedBy          *Person `json:"queuedBy"`
	PositionInSeconds int     `json:"positionInSeconds"`
	Playing           bool    `json:"playing"`
}

type Image struct {
	ID         string           `json:"id"`
	URL        string           `json:"url"`
	Thumbnails *ImageThumbnails `json:"thumbnails"`
	Sizes      *ImageSizes      `json:"sizes"`
}

type ImageSizes struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type ImageThumbnail struct {
	URL   string      `json:"url"`
	Sizes *ImageSizes `json:"sizes"`
}

type ImageThumbnails struct {
	S295  *ImageThumbnail `json:"s295"`
	M600  *ImageThumbnail `json:"m600"`
	B960  *ImageThumbnail `json:"b960"`
	W1200 *ImageThumbnail `json:"w1200"`
	F1920 *ImageThumbnail `json:"f1920"`
}

type PageInfo struct {
	StartCursor     string `json:"startCursor"`
	EndCursor       string `json:"endCursor"`
	HasNextPage     bool   `json:"hasNextPage"`
	HasPreviousPage bool   `json:"hasPreviousPage"`
}

type Person struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Player struct {
	ID            uuid.UUID      `json:"id"`
	Host          *Person        `json:"host"`
	CurrentTrack  *CurrentTrack  `json:"currentTrack,omitempty"`
	Queue         []*QueuedTrack `json:"queue"`
	TracksInQueue int            `json:"tracksInQueue"`
}

type QueuedTrack struct {
	Track    *Track    `json:"track"`
	QueuedBy *Person   `json:"queuedBy"`
	QueuedAt time.Time `json:"queuedAt"`
}

type QueuedTrackEdge struct {
	Cursor string       `json:"cursor"`
	Node   *QueuedTrack `json:"node"`
}

type QueuedTrackList struct {
	Edges    []*QueuedTrackEdge `json:"edges"`
	PageInfo *PageInfo          `json:"pageInfo"`
}

type Track struct {
	ID                uuid.UUID `json:"id"`
	Title             string    `json:"title"`
	DurationInSeconds int       `json:"durationInSeconds"`
	URL               string    `json:"url"`
}

type VoidBox struct {
	Value interface{} `json:"value"`
}
