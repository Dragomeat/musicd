package download_test

import (
	"musicd/internal/media"
	"musicd/internal/track/domain"
	"musicd/internal/track/download"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func TestWorkflow_Download(t *testing.T) {
	trackInfoFetcher := download.NewMockTrackInfoFetcher(t)
	trackInfoFetcher.EXPECT().
		Get(mock.Anything, "ZUkiNSvH6sA").
		Return(download.TrackInfo{
			ExternalID: "ZUkiNSvH6sA",
			Name:       "Iggy Pop - The Passenger",
			Files: map[download.MediaType]download.File{
				download.MediaTypeAudioMp4: {
					Url:     "https://www.youtube.com/watch?v=ZUkiNSvH6sA",
					Bitrate: 196000,
				},
			},
		}, nil).
		Once()
	uploader := media.NewMockUploader(t)
	uploader.EXPECT().
		Upload(mock.Anything, mock.Anything).
		Return(media.File{Id: "abobaa"}, nil).
		Once()

	workflow := download.NewWorkflow(trackInfoFetcher, uploader, domain.NullTracks{})
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	env.SetTestTimeout(10 * time.Second)
	workflow.Register(env)
	env.ExecuteWorkflow(download.WorkflowName, download.WorkflowInput{Url: "ZUkiNSvH6sA"})

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	var result string
	require.Error(t, env.GetWorkflowResult(&result))
	require.Empty(t, result)
}
