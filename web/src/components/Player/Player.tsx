import { useState, useEffect, useCallback, useRef, useMemo } from "react";
import { useQuery, useMutation } from "@apollo/client";
import { Box, Button, Spinner, Text, RangeInput, Layer } from "grommet";
import {
  Disc as DiscIcon,
  Menu as QueueIcon,
  Play as PlayIcon,
  Pause as PauseIcon,
  Rewind as PreviousIcon,
  FastForward as NextIcon,
  Volume as VolumeIcon,
  VolumeLow as VolumeLowIcon,
  VolumeMute as VolumeMuteIcon,
} from "grommet-icons";
import {
  GET_PLAYER,
  QUEUE_TRACK,
  REMOVE_TRACK_FROM_QUEUE,
  START_PLAYER,
  STOP_PLAYER,
  SEEK_TO,
  PREVIOUS_TRACK,
  NEXT_TRACK,
} from "./graphql";

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

// https://overreacted.io/making-setinterval-declarative-with-react-hooks/
const useInterval = (callback: Function, delay: number) => {
  const savedCallback = useRef<Function>();

  // Remember the latest callback.
  useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);

  // Set up the interval.
  useEffect(() => {
    function tick() {
      if (savedCallback.current) savedCallback.current();
    }
    if (delay !== null) {
      let id = setInterval(tick, delay);
      return () => clearInterval(id);
    }
  }, [delay]);
};

export const Player = ({ toggleQueueVisibility }: any) => {
  const { loading: loadingPlayer, data: { player } = {} } = useQuery(GET_PLAYER, {
    variables: { id: "1edebfee-dea1-6940-a35a-35daf92b7deb" },
  });

  const {
    id: playerId = "none",
    currentTrack: {
      track: { id = "none", title = "Unknown", url = "", durationInSeconds: defaultDurationInSeconds = "" } = {},
      positionInSeconds: defaultPositionInSeconds = 0,
      playing: defaultPlaying = false,
    } = {},
  } = player || {};

  const [startPlayer, { loading: startingPlayer }] = useMutation(START_PLAYER, {
    variables: {
      playerId,
    },
    onCompleted: () => setIsPlaying((prev) => !prev),
  });

  const [stopPlayer, { loading: stoppingPlayer }] = useMutation(STOP_PLAYER, {
    variables: {
      playerId,
    },
    onCompleted: () => setIsPlaying((prev) => !prev),
  });

  const [seekTo, { loading: seekingTo }] = useMutation(SEEK_TO);

  const [previousTrack, { loading: settingPreviousTrack }] = useMutation(PREVIOUS_TRACK);
  const [nextTrack, { data, loading: settingNextTrack }] = useMutation(NEXT_TRACK);

  const [isPlaying, setIsPlaying] = useState(false);
  const audioRef = useRef<HTMLAudioElement>(null);
  const playAnimationRef = useRef<number>();
  const togglePlayPause = () => {
    console.log(defaultPlaying, isPlaying);
    if (defaultPlaying !== isPlaying) {
      setIsPlaying((prev) => !prev);
      return;
    }

    if (isPlaying) {
      stopPlayer();
    } else {
      startPlayer();
    }
  };
  const progressBarRef = useRef<HTMLInputElement>(null);
  const handleProgressChange = () => {
    if (!audioRef.current) return;
    if (!progressBarRef.current) return;
    audioRef.current.currentTime = Number(progressBarRef.current.value);
  };
  const [timeProgress, setTimeProgress] = useState(0);
  const [duration, setDuration] = useState(0);

  useInterval(() => {
    if (!isPlaying) return;

    seekTo({
      variables: {
        playerId: playerId,
        positionInSeconds: Math.trunc(timeProgress),
      },
    });
  }, 5000);
  const onLoadedMetadata = () => {
    if (!audioRef.current) return;
    if (!progressBarRef.current) return;
    const seconds = audioRef.current.duration;
    setDuration(seconds);
    progressBarRef.current.max = String(seconds);
  };
  const repeat = useCallback(() => {
    if (!audioRef.current) return;
    if (!progressBarRef.current) return;
    const currentTime = audioRef.current.currentTime;
    setTimeProgress(currentTime);
    progressBarRef.current.value = String(currentTime);
    // progressBarRef.current.style.setProperty(
    //  "--range-progress",
    //  `${(Number(progressBarRef.current.value) / duration) * 100}%`
    // );

    playAnimationRef.current = requestAnimationFrame(repeat);
  }, [audioRef, progressBarRef, setTimeProgress]);
  useEffect(() => {
    if (!audioRef.current) return;
    if (isPlaying) {
      audioRef.current.play();
    } else {
      audioRef.current.pause();
    }
    playAnimationRef.current = requestAnimationFrame(repeat);
  }, [isPlaying, audioRef, repeat]);

  const [volume, setVolume] = useState(60);
  const [muteVolume, setMuteVolume] = useState(false);
  useEffect(() => {
    if (audioRef.current) {
      audioRef.current.volume = volume / 100;
      audioRef.current.muted = muteVolume;
    }
  }, [volume, muteVolume, audioRef]);

  useEffect(() => {
    if (isPlaying) {
      return;
    }
    if (!audioRef.current) return;
    if (!progressBarRef.current) return;
    audioRef.current.currentTime = defaultPositionInSeconds;
  }, [defaultPositionInSeconds]);
  useEffect(() => {
    setDuration(defaultDurationInSeconds);
  }, [defaultDurationInSeconds]);

  // Do not render if same track has different urls
  // TODO: what if url expired?
  const trackUrl = useMemo(() => url, [id]);

  const loading =
    loadingPlayer ||
    startingPlayer ||
    stoppingPlayer ||
    settingPreviousTrack ||
    settingNextTrack;

  return (
    <Box
      fill="horizontal"
      pad="small"
      gap="medium"
      direction="row"
      justify="between"
      background="light-1"
    >
      <Box direction="row" align="center" pad="small" gap="small" width="medium">
        <Box>
          <DiscIcon />
        </Box>
        <Box>
          <Text>{title}</Text>
          <Text>Unknown</Text>
        </Box>
      </Box>
      <Box pad="small" alignContent="center" justify="center" fill="horizontal">
        <Box direction="row" gap="small" alignSelf="center">
          <Button onClick={() => previousTrack({ variables: { playerId } })} disabled={loading}>
            <PreviousIcon />
          </Button>
          {loading ? (
            <Spinner />
          ) : (
            <Button onClick={togglePlayPause}>
              {isPlaying ? <PauseIcon /> : <PlayIcon />}
            </Button>
          )}
          <Button onClick={() => nextTrack({ variables: { playerId } })} disabled={loading}>
            <NextIcon />
          </Button>
        </Box>
        <Box direction="row" gap="small">
          <Text>{formatTime(timeProgress)}</Text>
          <RangeInput
            value={timeProgress}
            min={0}
            step={1}
            max={duration}
            onChange={handleProgressChange}
            ref={progressBarRef}
          />
          <Text>{formatTime(duration)}</Text>
        </Box>
        <audio
          autoPlay={isPlaying}
          src={trackUrl}
          onLoadedMetadata={onLoadedMetadata}
          onEnded={() => nextTrack({ variables: { playerId } })}
          ref={audioRef}
        />
      </Box>
      <Box
        pad="small"
        gap="medium"
        width="medium"
        direction="row"
        alignSelf="center"
        justify="end"
      >
        <Box>
          <Button onClick={toggleQueueVisibility}><QueueIcon /></Button>
        </Box>
        <Button onClick={() => setMuteVolume((prev) => !prev)}>
          {muteVolume || volume === 0 ? (
            <VolumeMuteIcon />
          ) : volume < 30 ? (
            <VolumeLowIcon />
          ) : (
            <VolumeIcon />
          )}
        </Button>

        <RangeInput
          min={0}
          step={1}
          max={100}
          value={volume}
          onChange={(e) => setVolume(Number(e.target.value))}
        />
      </Box>
    </Box>
  );
};
