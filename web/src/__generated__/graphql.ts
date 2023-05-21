/* eslint-disable */
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
  MilliTimestamp: any;
  Timestamp: any;
  UUID: any;
  Void: any;
};

export type CurrentTrack = {
  __typename?: 'CurrentTrack';
  playing: Scalars['Boolean'];
  positionInSeconds: Scalars['Int'];
  queuedBy: Person;
  track: Track;
};

export type Image = {
  __typename?: 'Image';
  id: Scalars['String'];
  sizes: ImageSizes;
  thumbnails: ImageThumbnails;
  url: Scalars['String'];
};

export type ImageSizes = {
  __typename?: 'ImageSizes';
  height: Scalars['Int'];
  width: Scalars['Int'];
};

export type ImageThumbnail = {
  __typename?: 'ImageThumbnail';
  sizes: ImageSizes;
  url: Scalars['String'];
};

export type ImageThumbnails = {
  __typename?: 'ImageThumbnails';
  b960: ImageThumbnail;
  f1920: ImageThumbnail;
  m600: ImageThumbnail;
  s295: ImageThumbnail;
  w1200: ImageThumbnail;
};

export type Mutation = {
  __typename?: 'Mutation';
  createPlayer: Player;
  moveTrackInQueue: Player;
  nextTrack: Player;
  previousTrack: Player;
  queueTrack: Player;
  removeTrackFromQueue: Player;
  seekTo: Player;
  startPlayer: Player;
  stopPlayer: Player;
};


export type MutationMoveTrackInQueueArgs = {
  playerId: Scalars['UUID'];
  position: Scalars['Int'];
  trackId: Scalars['UUID'];
};


export type MutationNextTrackArgs = {
  playerId: Scalars['UUID'];
};


export type MutationPreviousTrackArgs = {
  playerId: Scalars['UUID'];
};


export type MutationQueueTrackArgs = {
  playerId: Scalars['UUID'];
  trackId: Scalars['UUID'];
};


export type MutationRemoveTrackFromQueueArgs = {
  playerId: Scalars['UUID'];
  trackId: Scalars['UUID'];
};


export type MutationSeekToArgs = {
  playerId: Scalars['UUID'];
  positionInSeconds: Scalars['Int'];
};


export type MutationStartPlayerArgs = {
  playerId: Scalars['UUID'];
};


export type MutationStopPlayerArgs = {
  playerId: Scalars['UUID'];
};

export type PageInfo = {
  __typename?: 'PageInfo';
  endCursor: Scalars['String'];
  hasNextPage: Scalars['Boolean'];
  hasPreviousPage: Scalars['Boolean'];
  startCursor: Scalars['String'];
};

export type Person = {
  __typename?: 'Person';
  id: Scalars['UUID'];
  name: Scalars['String'];
};

export type Player = {
  __typename?: 'Player';
  currentTrack?: Maybe<CurrentTrack>;
  host: Person;
  id: Scalars['UUID'];
  queue: Array<QueuedTrack>;
  tracksInQueue: Scalars['Int'];
};


export type PlayerQueueArgs = {
  first?: Scalars['Int'];
};

export type Query = {
  __typename?: 'Query';
  player: Player;
  queue: QueuedTrackList;
};


export type QueryPlayerArgs = {
  playerId: Scalars['UUID'];
};


export type QueryQueueArgs = {
  after?: InputMaybe<Scalars['String']>;
  before?: InputMaybe<Scalars['String']>;
  first?: Scalars['Int'];
  playerId: Scalars['UUID'];
};

export type QueuedTrack = {
  __typename?: 'QueuedTrack';
  queuedAt: Scalars['Timestamp'];
  queuedBy: Person;
  track: Track;
};

export type QueuedTrackEdge = {
  __typename?: 'QueuedTrackEdge';
  cursor: Scalars['String'];
  node: QueuedTrack;
};

export type QueuedTrackList = {
  __typename?: 'QueuedTrackList';
  edges: Array<QueuedTrackEdge>;
  pageInfo: PageInfo;
};

export type Track = {
  __typename?: 'Track';
  durationInSeconds: Scalars['Int'];
  id: Scalars['UUID'];
  title: Scalars['String'];
  url: Scalars['String'];
};

export type VoidBox = {
  __typename?: 'VoidBox';
  value: Scalars['Void'];
};

export type PlayerFieldsFragment = { __typename?: 'Player', id: any, host: { __typename?: 'Person', id: any }, currentTrack?: { __typename?: 'CurrentTrack', positionInSeconds: number, playing: boolean, track: { __typename?: 'Track', id: any, title: string, url: string, durationInSeconds: number }, queuedBy: { __typename?: 'Person', id: any } } | null, queue: Array<{ __typename?: 'QueuedTrack', queuedAt: any, track: { __typename?: 'Track', id: any, title: string, durationInSeconds: number, url: string }, queuedBy: { __typename?: 'Person', id: any } }> } & { ' $fragmentName'?: 'PlayerFieldsFragment' };

export type GetPlayerQueryVariables = Exact<{
  id: Scalars['UUID'];
}>;


export type GetPlayerQuery = { __typename?: 'Query', player: (
    { __typename?: 'Player' }
    & { ' $fragmentRefs'?: { 'PlayerFieldsFragment': PlayerFieldsFragment } }
  ) };

export type StartPlayerMutationVariables = Exact<{
  playerId: Scalars['UUID'];
}>;


export type StartPlayerMutation = { __typename?: 'Mutation', startPlayer: (
    { __typename?: 'Player' }
    & { ' $fragmentRefs'?: { 'PlayerFieldsFragment': PlayerFieldsFragment } }
  ) };

export type StopPlayerMutationVariables = Exact<{
  playerId: Scalars['UUID'];
}>;


export type StopPlayerMutation = { __typename?: 'Mutation', stopPlayer: (
    { __typename?: 'Player' }
    & { ' $fragmentRefs'?: { 'PlayerFieldsFragment': PlayerFieldsFragment } }
  ) };

export type SeekToMutationVariables = Exact<{
  playerId: Scalars['UUID'];
  positionInSeconds: Scalars['Int'];
}>;


export type SeekToMutation = { __typename?: 'Mutation', seekTo: { __typename?: 'Player', id: any } };

export type PreviousTrackMutationVariables = Exact<{
  playerId: Scalars['UUID'];
}>;


export type PreviousTrackMutation = { __typename?: 'Mutation', previousTrack: (
    { __typename?: 'Player' }
    & { ' $fragmentRefs'?: { 'PlayerFieldsFragment': PlayerFieldsFragment } }
  ) };

export type NextTrackMutationVariables = Exact<{
  playerId: Scalars['UUID'];
}>;


export type NextTrackMutation = { __typename?: 'Mutation', nextTrack: (
    { __typename?: 'Player' }
    & { ' $fragmentRefs'?: { 'PlayerFieldsFragment': PlayerFieldsFragment } }
  ) };

export type QueueTrackMutationVariables = Exact<{
  playerId: Scalars['UUID'];
  trackId: Scalars['UUID'];
}>;


export type QueueTrackMutation = { __typename?: 'Mutation', queueTrack: (
    { __typename?: 'Player' }
    & { ' $fragmentRefs'?: { 'PlayerFieldsFragment': PlayerFieldsFragment } }
  ) };

export type RemoveTrackFromQueueMutationVariables = Exact<{
  playerId: Scalars['UUID'];
  trackId: Scalars['UUID'];
}>;


export type RemoveTrackFromQueueMutation = { __typename?: 'Mutation', removeTrackFromQueue: (
    { __typename?: 'Player' }
    & { ' $fragmentRefs'?: { 'PlayerFieldsFragment': PlayerFieldsFragment } }
  ) };

export const PlayerFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"PlayerFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Player"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"host"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"currentTrack"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"track"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"url"}},{"kind":"Field","name":{"kind":"Name","value":"durationInSeconds"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedBy"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"positionInSeconds"}},{"kind":"Field","name":{"kind":"Name","value":"playing"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queue"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"track"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"durationInSeconds"}},{"kind":"Field","name":{"kind":"Name","value":"url"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedBy"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedAt"}}]}}]}}]} as unknown as DocumentNode<PlayerFieldsFragment, unknown>;
export const GetPlayerDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetPlayer"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UUID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"player"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"playerId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"PlayerFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"PlayerFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Player"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"host"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"currentTrack"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"track"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"url"}},{"kind":"Field","name":{"kind":"Name","value":"durationInSeconds"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedBy"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"positionInSeconds"}},{"kind":"Field","name":{"kind":"Name","value":"playing"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queue"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"track"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"durationInSeconds"}},{"kind":"Field","name":{"kind":"Name","value":"url"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedBy"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedAt"}}]}}]}}]} as unknown as DocumentNode<GetPlayerQuery, GetPlayerQueryVariables>;
export const StartPlayerDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"StartPlayer"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"playerId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UUID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"startPlayer"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"playerId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"playerId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"PlayerFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"PlayerFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Player"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"host"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"currentTrack"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"track"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"url"}},{"kind":"Field","name":{"kind":"Name","value":"durationInSeconds"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedBy"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"positionInSeconds"}},{"kind":"Field","name":{"kind":"Name","value":"playing"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queue"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"track"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"durationInSeconds"}},{"kind":"Field","name":{"kind":"Name","value":"url"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedBy"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedAt"}}]}}]}}]} as unknown as DocumentNode<StartPlayerMutation, StartPlayerMutationVariables>;
export const StopPlayerDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"StopPlayer"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"playerId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UUID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"stopPlayer"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"playerId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"playerId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"PlayerFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"PlayerFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Player"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"host"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"currentTrack"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"track"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"url"}},{"kind":"Field","name":{"kind":"Name","value":"durationInSeconds"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedBy"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"positionInSeconds"}},{"kind":"Field","name":{"kind":"Name","value":"playing"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queue"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"track"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"durationInSeconds"}},{"kind":"Field","name":{"kind":"Name","value":"url"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedBy"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedAt"}}]}}]}}]} as unknown as DocumentNode<StopPlayerMutation, StopPlayerMutationVariables>;
export const SeekToDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"SeekTo"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"playerId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UUID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"positionInSeconds"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"seekTo"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"playerId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"playerId"}}},{"kind":"Argument","name":{"kind":"Name","value":"positionInSeconds"},"value":{"kind":"Variable","name":{"kind":"Name","value":"positionInSeconds"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]} as unknown as DocumentNode<SeekToMutation, SeekToMutationVariables>;
export const PreviousTrackDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"PreviousTrack"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"playerId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UUID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"previousTrack"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"playerId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"playerId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"PlayerFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"PlayerFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Player"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"host"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"currentTrack"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"track"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"url"}},{"kind":"Field","name":{"kind":"Name","value":"durationInSeconds"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedBy"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"positionInSeconds"}},{"kind":"Field","name":{"kind":"Name","value":"playing"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queue"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"track"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"durationInSeconds"}},{"kind":"Field","name":{"kind":"Name","value":"url"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedBy"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedAt"}}]}}]}}]} as unknown as DocumentNode<PreviousTrackMutation, PreviousTrackMutationVariables>;
export const NextTrackDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"NextTrack"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"playerId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UUID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"nextTrack"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"playerId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"playerId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"PlayerFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"PlayerFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Player"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"host"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"currentTrack"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"track"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"url"}},{"kind":"Field","name":{"kind":"Name","value":"durationInSeconds"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedBy"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"positionInSeconds"}},{"kind":"Field","name":{"kind":"Name","value":"playing"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queue"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"track"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"durationInSeconds"}},{"kind":"Field","name":{"kind":"Name","value":"url"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedBy"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedAt"}}]}}]}}]} as unknown as DocumentNode<NextTrackMutation, NextTrackMutationVariables>;
export const QueueTrackDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"QueueTrack"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"playerId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UUID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"trackId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UUID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"queueTrack"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"playerId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"playerId"}}},{"kind":"Argument","name":{"kind":"Name","value":"trackId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"trackId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"PlayerFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"PlayerFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Player"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"host"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"currentTrack"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"track"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"url"}},{"kind":"Field","name":{"kind":"Name","value":"durationInSeconds"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedBy"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"positionInSeconds"}},{"kind":"Field","name":{"kind":"Name","value":"playing"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queue"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"track"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"durationInSeconds"}},{"kind":"Field","name":{"kind":"Name","value":"url"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedBy"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedAt"}}]}}]}}]} as unknown as DocumentNode<QueueTrackMutation, QueueTrackMutationVariables>;
export const RemoveTrackFromQueueDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"RemoveTrackFromQueue"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"playerId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UUID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"trackId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UUID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"removeTrackFromQueue"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"playerId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"playerId"}}},{"kind":"Argument","name":{"kind":"Name","value":"trackId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"trackId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"PlayerFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"PlayerFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Player"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"host"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"currentTrack"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"track"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"url"}},{"kind":"Field","name":{"kind":"Name","value":"durationInSeconds"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedBy"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"positionInSeconds"}},{"kind":"Field","name":{"kind":"Name","value":"playing"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queue"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"track"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"durationInSeconds"}},{"kind":"Field","name":{"kind":"Name","value":"url"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedBy"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}},{"kind":"Field","name":{"kind":"Name","value":"queuedAt"}}]}}]}}]} as unknown as DocumentNode<RemoveTrackFromQueueMutation, RemoveTrackFromQueueMutationVariables>;