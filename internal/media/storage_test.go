package media_test

import (
	"context"
	"musicd/internal/media"
	"os"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStorage__Upload(t *testing.T) {
	uploadId := uuid.Must(uuid.NewV4()).String()

	audio, err := os.Open("testdata/iggy-pop-the-passenger.mp3")
	require.NoError(t, err)

	defer audio.Close()
	driver := media.NewMockDriver(t)
	driver.EXPECT().
		Put(
			context.Background(),
			mock.MatchedBy(func(upload media.Upload) bool {
				return upload.Id == uploadId &&
					upload.FileName == "Iggy Pop - The Passenger.mp3" &&
					upload.ContentType == "audio/mpeg" &&
					upload.Size == 4547115
			}),
		).
		Return(
			media.File{
				Id:          uploadId,
				Sha256:      "06dd894d08ad73ada644d2231dfc652538d86852a01c789779419f4600fa3956",
				FileName:    "Iggy Pop - The Passenger.mp3",
				ContentType: "audio/mpeg",
				Size:        4547115,
			},
			nil,
		).
		Once()
	driver.EXPECT().
		Move(
			context.Background(),
			uploadId,
			"06dd894d08ad73ada644d2231dfc652538d86852a01c789779419f4600fa3956",
		).
		Return(nil).
		Once()

	storage := media.NewStorage(driver)

	file, err := storage.Upload(
		context.Background(),
		media.Upload{
			Id:          uploadId,
			File:        audio,
			FileName:    "Iggy Pop - The Passenger.mp3",
			ContentType: "audio/mpeg",
			Size:        4547115,
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, "06dd894d08ad73ada644d2231dfc652538d86852a01c789779419f4600fa3956", file.Id)
	assert.Equal(t, "Iggy Pop - The Passenger.mp3", file.FileName)
	assert.Equal(t, "audio/mpeg", file.ContentType)
	assert.Equal(t, 4547115, file.Size)
}

func TestStorage__GetLink(t *testing.T) {
	driver := media.NewMockDriver(t)
	driver.EXPECT().
		GetMetadata(context.Background(), "06dd894d08ad73ada644d2231dfc652538d86852a01c789779419f4600fa3956").
		Return(
			media.File{
				Id:          "06dd894d08ad73ada644d2231dfc652538d86852a01c789779419f4600fa3956",
				FileName:    "Iggy Pop - The Passenger.mp3",
				ContentType: "audio/mpeg",
				Size:        4547115,
			},
			nil,
		).
		Once()

	storage := media.NewStorage(driver)

	link, err := storage.GetLink(context.Background(), "06dd894d08ad73ada644d2231dfc652538d86852a01c789779419f4600fa3956")
	assert.NoError(t, err)
	assert.Equal(t, "http://cdn.music.local/06dd894d08ad73ada644d2231dfc652538d86852a01c789779419f4600fa3956.mp3", link)
}

func TestStorage__GetLink__NonExistingFile(t *testing.T) {
	driver := media.NewMockDriver(t)
	driver.EXPECT().GetMetadata(context.Background(), "non-existing-file").Return(media.File{}, media.ErrFileNotFound).Once()

	storage := media.NewStorage(driver)

	link, err := storage.GetLink(context.Background(), "non-existing-file")
	assert.ErrorIs(t, err, media.ErrFileNotFound)
	assert.Empty(t, link)
}
