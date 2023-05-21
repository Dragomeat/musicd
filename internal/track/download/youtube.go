package download

import (
	"net/http"

	"github.com/kkdai/youtube/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/net/context"
)

type YoutubeDownloader struct{}

func NewYoutubeDownloader() *YoutubeDownloader {
	return &YoutubeDownloader{}
}

func (d *YoutubeDownloader) Get(ctx context.Context, url string) (TrackInfo, error) {
	client := youtube.Client{
		HTTPClient: &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)},
		ChunkSize:  youtube.Size1Mb,
	}
	video, err := client.GetVideoContext(ctx, url)
	if err != nil {
		return TrackInfo{}, err
	}

	format := video.Formats.Type("audio/mp4")[0]
	sUrl, err := client.GetStreamURLContext(ctx, video, &format)
	if err != nil {
		return TrackInfo{}, err
	}
	return TrackInfo{
		ExternalID:        video.ID,
		Name:              video.Title,
		DurationInSeconds: int(video.Duration.Seconds()),
		Files:             map[MediaType]File{MediaTypeAudioMp4: {Url: sUrl, Bitrate: format.Bitrate}},
	}, nil
}
