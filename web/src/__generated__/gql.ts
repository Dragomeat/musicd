/* eslint-disable */
import * as types from './graphql';
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 */
const documents = {
    "\n  fragment PlayerFields on Player {\n    id\n    host {\n      id\n    }\n    currentTrack {\n      track {\n        id\n        title\n        url\n        durationInSeconds\n      }\n      queuedBy {\n        id\n      }\n      positionInSeconds\n      playing\n    }\n    queue {\n      track {\n        id\n        title\n        durationInSeconds\n        url\n      }\n      queuedBy {\n        id\n      }\n      queuedAt\n    }\n  }\n": types.PlayerFieldsFragmentDoc,
    "\n  query GetPlayer($id: UUID!) {\n    player(playerId: $id) {\n      ...PlayerFields\n    }\n  }\n  \n": types.GetPlayerDocument,
    "\n  mutation StartPlayer($playerId: UUID!) {\n    startPlayer(playerId: $playerId) {\n      ...PlayerFields\n    }\n  }\n  \n": types.StartPlayerDocument,
    "\n  mutation StopPlayer($playerId: UUID!) {\n    stopPlayer(playerId: $playerId) {\n      ...PlayerFields\n    }\n  }\n  \n": types.StopPlayerDocument,
    "\n  mutation SeekTo($playerId: UUID!, $positionInSeconds: Int!) {\n    seekTo(playerId: $playerId, positionInSeconds: $positionInSeconds) {\n      id\n    }\n  }\n": types.SeekToDocument,
    "\n  mutation PreviousTrack($playerId: UUID!) {\n    previousTrack(playerId: $playerId) {\n      ...PlayerFields\n    }\n  }\n  \n": types.PreviousTrackDocument,
    "\n  mutation NextTrack($playerId: UUID!) {\n    nextTrack(playerId: $playerId) {\n      ...PlayerFields\n    }\n  }\n  \n": types.NextTrackDocument,
    "\n  mutation QueueTrack($playerId: UUID!, $trackId: UUID!) {\n    queueTrack(playerId: $playerId, trackId: $trackId) {\n      ...PlayerFields\n    }\n  }\n  \n": types.QueueTrackDocument,
    "\n  mutation RemoveTrackFromQueue($playerId: UUID!, $trackId: UUID!) {\n    removeTrackFromQueue(playerId: $playerId, trackId: $trackId) {\n      ...PlayerFields\n    }\n  }\n  \n": types.RemoveTrackFromQueueDocument,
};

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = gql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function gql(source: string): unknown;

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  fragment PlayerFields on Player {\n    id\n    host {\n      id\n    }\n    currentTrack {\n      track {\n        id\n        title\n        url\n        durationInSeconds\n      }\n      queuedBy {\n        id\n      }\n      positionInSeconds\n      playing\n    }\n    queue {\n      track {\n        id\n        title\n        durationInSeconds\n        url\n      }\n      queuedBy {\n        id\n      }\n      queuedAt\n    }\n  }\n"): (typeof documents)["\n  fragment PlayerFields on Player {\n    id\n    host {\n      id\n    }\n    currentTrack {\n      track {\n        id\n        title\n        url\n        durationInSeconds\n      }\n      queuedBy {\n        id\n      }\n      positionInSeconds\n      playing\n    }\n    queue {\n      track {\n        id\n        title\n        durationInSeconds\n        url\n      }\n      queuedBy {\n        id\n      }\n      queuedAt\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query GetPlayer($id: UUID!) {\n    player(playerId: $id) {\n      ...PlayerFields\n    }\n  }\n  \n"): (typeof documents)["\n  query GetPlayer($id: UUID!) {\n    player(playerId: $id) {\n      ...PlayerFields\n    }\n  }\n  \n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation StartPlayer($playerId: UUID!) {\n    startPlayer(playerId: $playerId) {\n      ...PlayerFields\n    }\n  }\n  \n"): (typeof documents)["\n  mutation StartPlayer($playerId: UUID!) {\n    startPlayer(playerId: $playerId) {\n      ...PlayerFields\n    }\n  }\n  \n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation StopPlayer($playerId: UUID!) {\n    stopPlayer(playerId: $playerId) {\n      ...PlayerFields\n    }\n  }\n  \n"): (typeof documents)["\n  mutation StopPlayer($playerId: UUID!) {\n    stopPlayer(playerId: $playerId) {\n      ...PlayerFields\n    }\n  }\n  \n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation SeekTo($playerId: UUID!, $positionInSeconds: Int!) {\n    seekTo(playerId: $playerId, positionInSeconds: $positionInSeconds) {\n      id\n    }\n  }\n"): (typeof documents)["\n  mutation SeekTo($playerId: UUID!, $positionInSeconds: Int!) {\n    seekTo(playerId: $playerId, positionInSeconds: $positionInSeconds) {\n      id\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation PreviousTrack($playerId: UUID!) {\n    previousTrack(playerId: $playerId) {\n      ...PlayerFields\n    }\n  }\n  \n"): (typeof documents)["\n  mutation PreviousTrack($playerId: UUID!) {\n    previousTrack(playerId: $playerId) {\n      ...PlayerFields\n    }\n  }\n  \n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation NextTrack($playerId: UUID!) {\n    nextTrack(playerId: $playerId) {\n      ...PlayerFields\n    }\n  }\n  \n"): (typeof documents)["\n  mutation NextTrack($playerId: UUID!) {\n    nextTrack(playerId: $playerId) {\n      ...PlayerFields\n    }\n  }\n  \n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation QueueTrack($playerId: UUID!, $trackId: UUID!) {\n    queueTrack(playerId: $playerId, trackId: $trackId) {\n      ...PlayerFields\n    }\n  }\n  \n"): (typeof documents)["\n  mutation QueueTrack($playerId: UUID!, $trackId: UUID!) {\n    queueTrack(playerId: $playerId, trackId: $trackId) {\n      ...PlayerFields\n    }\n  }\n  \n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation RemoveTrackFromQueue($playerId: UUID!, $trackId: UUID!) {\n    removeTrackFromQueue(playerId: $playerId, trackId: $trackId) {\n      ...PlayerFields\n    }\n  }\n  \n"): (typeof documents)["\n  mutation RemoveTrackFromQueue($playerId: UUID!, $trackId: UUID!) {\n    removeTrackFromQueue(playerId: $playerId, trackId: $trackId) {\n      ...PlayerFields\n    }\n  }\n  \n"];

export function gql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;