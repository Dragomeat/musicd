package download

import (
	"musicd/internal/media"
	"musicd/internal/track/domain"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"golang.org/x/net/context"
)

const (
	WorkflowName              = "track.download"
	activityNameDownloadAudio = "track.downloadAudio"
	activityNameCreate        = "track.create"
)

type MediaType string

const (
	MediaTypeAudioMp4 MediaType = "audio/mp4"
)

type File struct {
	Url     string
	Bitrate int
}

type TrackInfo struct {
	ExternalID        string
	Name              string
	DurationInSeconds int
	Files             map[MediaType]File
}

//go:generate go run github.com/vektra/mockery/v2 --name TrackInfoFetcher
type TrackInfoFetcher interface {
	Get(ctx context.Context, url string) (TrackInfo, error)
}

type WorkflowInput struct {
	Url string
}

type Workflow struct {
	trackInfoFetcher TrackInfoFetcher
	uploader         media.Uploader
	tracks           domain.Tracks
}

func NewWorkflow(
	trackInfoFetcher TrackInfoFetcher,
	uploader media.Uploader,
	tracks domain.Tracks,
) *Workflow {
	return &Workflow{
		trackInfoFetcher: trackInfoFetcher,
		uploader:         uploader,
		tracks:           tracks,
	}
}

func (w *Workflow) Register(registry worker.Registry) {
	registry.RegisterWorkflowWithOptions(
		w.Download,
		workflow.RegisterOptions{Name: WorkflowName},
	)

	registry.RegisterActivityWithOptions(
		w.downloadTrack,
		activity.RegisterOptions{Name: activityNameDownloadAudio},
	)
	registry.RegisterActivityWithOptions(
		w.createTrack,
		activity.RegisterOptions{Name: activityNameCreate},
	)
}

func (w *Workflow) Download(ctx workflow.Context, input WorkflowInput) error {
	ctx = workflow.WithActivityOptions(
		ctx,
		workflow.ActivityOptions{
			StartToCloseTimeout: 5 * time.Second,
		},
	)

	var downloadTrackOutput downloadTrackOutput
	err := workflow.
		ExecuteActivity(
			workflow.WithActivityOptions(
				ctx,
				workflow.ActivityOptions{
					StartToCloseTimeout: 5 * time.Minute,
					HeartbeatTimeout:    2 * time.Second,
				},
			),
			activityNameDownloadAudio,
			downloadTrackInput{Url: input.Url},
		).
		Get(ctx, &downloadTrackOutput)
	if err != nil {
		return err
	}

	workflow.GetLogger(ctx).Info("Track info fetched", "trackInfo", downloadTrackOutput.TrackInfo)
	workflow.GetLogger(ctx).Info("Audio downloaded", "fileId", downloadTrackOutput.FileId)

	var createTrackOutput createTrackOutput
	err = workflow.
		ExecuteActivity(
			ctx,
			activityNameCreate,
			createTrackInput{Track: downloadTrackOutput.TrackInfo, FileId: downloadTrackOutput.FileId},
		).
		Get(ctx, &createTrackOutput)
	if err != nil {
		return err
	}
	workflow.GetLogger(ctx).Info("Track created", "trackId", createTrackOutput.Track.ID)

	return nil
}

type downloadTrackInput struct {
	Url string
}
type downloadTrackOutput struct {
	FileId    string
	TrackInfo TrackInfo
}

func (w *Workflow) downloadTrack(ctx context.Context, input downloadTrackInput) (downloadTrackOutput, error) {
	trackInfo, err := w.trackInfoFetcher.Get(ctx, input.Url)
	if err != nil {
		return downloadTrackOutput{}, err
	}

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, trackInfo.Files[MediaTypeAudioMp4].Url, nil)
	if err != nil {
		return downloadTrackOutput{}, err
	}
	res, err := client.Do(req)
	if err != nil {
		return downloadTrackOutput{}, err
	}
	defer res.Body.Close()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Second):
				activity.RecordHeartbeat(ctx)
			}
		}
	}()

	file, err := w.uploader.Upload(ctx, media.Upload{
		Id:          uuid.Must(uuid.NewV4()).String(),
		File:        res.Body,
		FileName:    trackInfo.Name,
		ContentType: res.Header.Get("Content-Type"),
		Size:        int(res.ContentLength),
	})
	if err != nil {
		return downloadTrackOutput{}, err
	}

	return downloadTrackOutput{FileId: file.Id, TrackInfo: trackInfo}, nil
}

type createTrackInput struct {
	Track  TrackInfo
	FileId string
}
type createTrackOutput struct {
	Track domain.Track
}

func (w *Workflow) createTrack(ctx context.Context, input createTrackInput) (createTrackOutput, error) {
	track, err := domain.NewTrack(
		input.Track.Name,
		input.Track.DurationInSeconds,
		[]domain.Audio{
			{
				ExternalId:     input.Track.ExternalID,
				ExternalSource: domain.ExternalSourceYoutube,
				Sha256:         input.FileId,
				Bitrate:        input.Track.Files[MediaTypeAudioMp4].Bitrate,
				Type:           domain.MIMETypeAudioMp4,
			},
		},
		uuid.Nil,
	)
	if err != nil {
		return createTrackOutput{}, err
	}

	err = w.tracks.Create(ctx, track)
	if err != nil {
		return createTrackOutput{}, err
	}

	return createTrackOutput{Track: track}, nil
}
