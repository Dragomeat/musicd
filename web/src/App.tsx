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

const Track = ({ track, addTrackToQueue, ...props }: any) => (
  <Box direction="row" align="center" pad={{ top: "small", bottom: "small" }} gap="small" {...props}>
    <Box direction="row" gap="small" align="center" fill="horizontal">
      <DiscIcon />
      <Box>
        <Text>{track.title}</Text>
        <Text>Unknown</Text>
      </Box>
    </Box>
    <Box width="medium">
      <Text>Unknown</Text>
    </Box>
    <Box direction="row" gap="medium" justify="end">
      <Text>{formatTime(track.durationInSeconds)}</Text>
      <Button onClick={() => addTrackToQueue(track.id)}>
        <AddIcon />
      </Button>
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
            <Box overflow={{ vertical: "scroll" }}>
              <List data={data?.tracks?.edges || []} step={10} border={false} onMore={() => fetchMore({
                variables: {
                  after: data?.tracks?.pageInfo?.endCursor,
                },
              })}>
                {(edge: any) => (<Track key={edge.node.id} track={edge.node} addTrackToQueue={(trackId: string) => queueTrack({ variables: { playerId: player.id, trackId } })} />)}
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
              <Box background="light-2" fill="horizontal" height="medium" pad="medium">
                {player.queue.map((queuedTrack: any, i: number) => (
                  <Box direction="row" align="center" pad={{ top: "small", bottom: "small" }} gap="small" key={`qt-${queuedTrack.track.id}`}>
                    <Box direction="row" gap="small" align="center" fill="horizontal">
                      <DiscIcon />
                      <Box>
                        <Text>{queuedTrack.track.title}</Text>
                        <Text>Unknown</Text>
                      </Box>
                    </Box>
                    <Box width="medium">
                      <Text>Unknown</Text>
                    </Box>
                    <Box direction="row" gap="medium" justify="end">
                      <Text>{formatTime(queuedTrack.track.durationInSeconds)}</Text>
                      <Button
                        onClick={() =>
                          removeTrackFromQueue({
                            variables: {
                              playerId: player.id,
                              trackId: queuedTrack.track.id,
                            },
                          })
                        }
                      >
                        <RemoveIcon />
                      </Button>
                    </Box>
                  </Box>
                ))}
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
