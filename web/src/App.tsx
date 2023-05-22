import { useState, useRef } from "react";
import { useQuery, useMutation, gql } from "@apollo/client";
import {
  Grommet,
  Page,
  PageContent,
  Heading,
  Box,
  List,
  Grid,
  Main,
  Footer,
  Text,
  Layer,
  Button,
  InfiniteScroll,
} from "grommet";
import {
  Disc as DiscIcon,
  Add as AddIcon,
  Trash as RemoveIcon,
} from "grommet-icons";
import { Player } from "./components/Player/Player";
import {
  GET_PLAYER,
  QUEUE_TRACK,
  REMOVE_TRACK_FROM_QUEUE,
} from "./components/Player/graphql";

export const PAGAINATE_TRACKS = gql`
  query PaginateTracks($first: Int!, $after: String) {
     tracks(after: $after, first: $first) {
        edges {
           cursor
           node {
             id
             title
             durationInSeconds
          }
        }
        pageInfo {
          startCursor
          endCursor
          hasNextPage
          hasPreviousPage
        }
    }
  }
`;

const formatTime = (time: number) => {
  if (time && !isNaN(time)) {
    const minutes = Math.floor(time / 60);
    const formatMinutes = minutes < 10 ? `0${minutes}` : `${minutes}`;
    const seconds = Math.floor(time % 60);
    const formatSeconds = seconds < 10 ? `0${seconds}` : `${seconds}`;
    return `${formatMinutes}:${formatSeconds}`;
  }
  return "00:00";
};

const Track = ({ track, inQueue, addTrackToQueue, removeTrackFromQueue, ...props }: any) => (
  <Box direction="row" align="center" pad={{ top: "small", bottom: "small" }} justify="between" gap="medium" {...props}>
    <Box direction="row" gap="small" align="center" width="large">
      <DiscIcon />
      <Box>
        <Text truncate="tip">{track.title}</Text>
        <Text>Unknown</Text>
      </Box>
    </Box>
    <Box fill="horizontal" align="center">
      <Text>Unknown</Text>
    </Box>
    <Box direction="row" width="small" gap="medium" justify="end">
      <Text>{formatTime(track.durationInSeconds)}</Text>
      {
        inQueue
          ? (
            <Button onClick={() => removeTrackFromQueue(track.id)}>
              <RemoveIcon />
            </Button>
          )
          : (
            <Button onClick={() => addTrackToQueue(track.id)}>
              <AddIcon />
            </Button>
          )
      }
    </Box>
  </Box>
);

function App() {
  const { loading: loadingPlayer, data: { player } = {} } = useQuery(GET_PLAYER, {
    variables: { id: "1edebfee-dea1-6940-a35a-35daf92b7deb" },
  });
  const { loading, error, data, fetchMore } = useQuery(PAGAINATE_TRACKS, {
    variables: { first: 10 },
  });

  const [queueTrack] = useMutation(QUEUE_TRACK);
  const [removeTrackFromQueue] = useMutation(REMOVE_TRACK_FROM_QUEUE);

  const [isQueueOpen, setIsQueueOpen] = useState(false);

  return (
    <Grommet full>
      <Page kind="full">
        <PageContent>
          <Main>
            <Heading>Tracks</Heading>
            <Box overflow={{ vertical: "scroll" }} pad={{ bottom: "xlarge" }}>
              <List data={data?.tracks?.edges || []} pad="none" step={10} border={false} onMore={() => fetchMore({
                variables: {
                  after: data?.tracks?.pageInfo?.endCursor,
                },
              })}>
                {(edge: any) => (<Track key={edge.node.id} track={edge.node} inQueue={false} addTrackToQueue={(trackId: string) => queueTrack({ variables: { playerId: player.id, trackId } })} />)}
              </List>
            </Box>
          </Main>
          <Layer
            position="bottom"
            full="horizontal"
            modal={false}
            plain
            responsive={false}
            onEsc={() => setIsQueueOpen(false)}
            onClickOutside={() => setIsQueueOpen(false)}
          >
            {isQueueOpen && (
              <Box background="light-2" fill="horizontal" height="medium" pad={{ top: "small", left: "medium", right: "medium", bottom: "medium" }}>
                <Heading margin={{ bottom: "small" }} level={3}>Queue</Heading>
                <Box overflow={{ vertical: "scroll" }} fill="vertical">
                  {player?.currentTrack && (
                    <Box flex={false}>
                      <Heading margin="none" level={4}>Now playing</Heading>
                      <Box flex={false}>
                        <Track
                          track={player.currentTrack.track}
                          inQueue={true}
                          addTrackToQueue={(trackId: string) => queueTrack({ variables: { playerId: player.id, trackId } })}
                          removeTrackFromQueue={(trackId: string) => removeTrackFromQueue({ variables: { playerId: player.id, trackId } })}
                        />
                      </Box>
                    </Box>
                  )}
                  <Box flex={false}>
                    <Heading margin="none" level={4}>Next</Heading>
                    <Box flex={false}>
                      {
                        player.queue?.length > 0
                          ? (
                            <List data={player.queue} pad="none" step={5} border={false}>
                              {(queuedTrack: any) => (
                                <Track
                                  key={queuedTrack.track.id}
                                  inQueue={true}
                                  track={queuedTrack.track}
                                  addTrackToQueue={(trackId: string) => queueTrack({ variables: { playerId: player.id, trackId } })}
                                  removeTrackFromQueue={(trackId: string) => removeTrackFromQueue({ variables: { playerId: player.id, trackId } })}
                                />
                              )}
                            </List>
                          ) : (
                            <Text>No tracks in queue</Text>
                          )
                      }
                    </Box>
                  </Box>
                </Box>
              </Box>
            )}
            <Player toggleQueueVisibility={() => setIsQueueOpen((prev) => !prev)} />
          </Layer>
        </PageContent>
      </Page>
    </Grommet>
  );
}

export default App;
