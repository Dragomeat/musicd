import { gql, TypedDocumentNode } from "@apollo/client";

export const PLAYER_FIELDS_FRAGMENT = gql`
  fragment PlayerFields on Player {
    id
    host {
      id
    }
    currentTrack {
      track {
        id
        title
        url
        durationInSeconds
      }
      queuedBy {
        id
      }
      positionInSeconds
      playing
    }
    queue {
      track {
        id
        title
        durationInSeconds
        url
      }
      queuedBy {
        id
      }
      queuedAt
    }
  }
`;

export const GET_PLAYER = gql`
  query GetPlayer($id: UUID!) {
    player(playerId: $id) {
      ...PlayerFields
    }
  }
  ${PLAYER_FIELDS_FRAGMENT}
`;

export const START_PLAYER = gql`
  mutation StartPlayer($playerId: UUID!) {
    startPlayer(playerId: $playerId) {
      ...PlayerFields
    }
  }
  ${PLAYER_FIELDS_FRAGMENT}
`;

export const STOP_PLAYER = gql`
  mutation StopPlayer($playerId: UUID!) {
    stopPlayer(playerId: $playerId) {
      ...PlayerFields
    }
  }
  ${PLAYER_FIELDS_FRAGMENT}
`;

export const SEEK_TO = gql`
  mutation SeekTo($playerId: UUID!, $positionInSeconds: Int!) {
    seekTo(playerId: $playerId, positionInSeconds: $positionInSeconds) {
      ...PlayerFields
    }
  }
  ${PLAYER_FIELDS_FRAGMENT}
`;

export const PREVIOUS_TRACK = gql`
  mutation PreviousTrack($playerId: UUID!) {
    previousTrack(playerId: $playerId) {
      ...PlayerFields
    }
  }
  ${PLAYER_FIELDS_FRAGMENT}
`;

export const NEXT_TRACK = gql`
  mutation NextTrack($playerId: UUID!) {
    nextTrack(playerId: $playerId) {
      ...PlayerFields
    }
  }
  ${PLAYER_FIELDS_FRAGMENT}
`;

export const QUEUE_TRACK = gql`
  mutation QueueTrack($playerId: UUID!, $trackId: UUID!) {
    queueTrack(playerId: $playerId, trackId: $trackId) {
      ...PlayerFields
    }
  }
  ${PLAYER_FIELDS_FRAGMENT}
`;

export const REMOVE_TRACK_FROM_QUEUE = gql`
  mutation RemoveTrackFromQueue($playerId: UUID!, $trackId: UUID!) {
    removeTrackFromQueue(playerId: $playerId, trackId: $trackId) {
      ...PlayerFields
    }
  }
  ${PLAYER_FIELDS_FRAGMENT}
`;
