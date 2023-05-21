import { useState } from "react";
import { useQuery, useMutation, gql } from "@apollo/client";
import {
  Grommet,
  Page,
  PageContent,
  Box,
  Button,
  Text,
  TextInput,
} from "grommet";
import { Trash as RemoveIcon } from "grommet-icons";
import { Player } from "./components/Player/Player";
import {
  GET_PLAYER,
  QUEUE_TRACK,
  REMOVE_TRACK_FROM_QUEUE,
} from "./components/Player/graphql";

function App() {
  const { loading, error, data } = useQuery(GET_PLAYER, {
    variables: { id: "1edebfee-dea1-6940-a35a-35daf92b7deb" },
  });

  const [trackId, setTrackId] = useState("");
  const [queueTrack] = useMutation(QUEUE_TRACK, {
    onCompleted: () => setTrackId(""),
  });
  const [removeTrackFromQueue] = useMutation(REMOVE_TRACK_FROM_QUEUE);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error : {error.message}</p>;

  return (
    <Grommet full>
      <Page>
        <PageContent>
          <Box direction="row">
            <TextInput
              placeholder="Search"
              value={trackId}
              onChange={(e) => setTrackId(e.target.value)}
            />
            <Button
              primary
              onClick={() =>
                queueTrack({ variables: { playerId: data.player.id, trackId } })
              }
              label="Add"
            />
          </Box>
          <Box>
            {data.player.queue.map((queuedTrack: any, i: number) => (
              <Box
                key={`qt${queuedTrack.track.id}`}
                direction="row"
                gap="small"
              >
                <Text>
                  {i + 1}. {queuedTrack.track.title}
                </Text>
                <Button
                  onClick={() =>
                    removeTrackFromQueue({
                      variables: {
                        playerId: data.player.id,
                        trackId: queuedTrack.track.id,
                      },
                    })
                  }
                >
                  <RemoveIcon />
                </Button>
              </Box>
            ))}
          </Box>
          <Player {...data.player} />
        </PageContent>
      </Page>
    </Grommet>
  );
}

export default App;
