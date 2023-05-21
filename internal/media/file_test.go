package media_test

import (
	"context"
	"musicd/internal/media"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileDriver__Put(t *testing.T) {
	driver := media.NewFileDriver("/tmp")

	audio, err := os.Open("testdata/iggy-pop-the-passenger.mp3")
	assert.NoError(t, err)
	defer audio.Close()

	err = driver.Put(
		context.Background(),
		media.File{
			Id:          "06dd894d08ad73ada644d2231dfc652538d86852a01c789779419f4600fa3956",
			FileName:    "Iggy Pop - The Passenger.mp3",
			ContentType: "audio/mpeg",
			Size:        4547115,
		}, audio,
	)
	assert.NoError(t, err)
}

func TestFileDriver__GetMetadata(t *testing.T) {
	driver := media.NewFileDriver("/tmp")

	audio, err := os.Open("testdata/iggy-pop-the-passenger.mp3")
	assert.NoError(t, err)
	defer audio.Close()

	err = driver.Put(
		context.Background(),
		media.File{
			Id:          "06dd894d08ad73ada644d2231dfc652538d86852a01c789779419f4600fa3956",
			FileName:    "Iggy Pop - The Passenger.mp3",
			ContentType: "audio/mpeg",
			Size:        4547115,
		}, audio,
	)
	assert.NoError(t, err)

	file, err := driver.GetMetadata(context.Background(), "06dd894d08ad73ada644d2231dfc652538d86852a01c789779419f4600fa3956")
	assert.NoError(t, err)

	assert.Equal(t, "06dd894d08ad73ada644d2231dfc652538d86852a01c789779419f4600fa3956", file.Id)
	assert.Equal(t, "Iggy Pop - The Passenger.mp3", file.FileName)
	assert.Equal(t, "audio/mpeg", file.ContentType)
	assert.Equal(t, 4547115, file.Size)
}
