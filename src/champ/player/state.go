package player

import "fmt"

type State int

const (
	STATE_STOPPED   State = 0
	STATE_PLAYING         = 1
	STATE_PAUSED          = 2
	STATE_BUFFERING       = 3
	STATE_SEEKING         = 4
)

type PlayState struct {
	State    State
	Volume   int32
	Position int32
	Duration int32
}

func (s *PlayState) String() string {
	return fmt.Sprintf(
		"state:%d volume:%d position:%d/%d",
		s.State, s.Volume, s.Position, s.Duration,
	)
}

type CoreEventType int

const (
	CoreReady           CoreEventType = 0
	CorePause           CoreEventType = 1
	CoreSeek            CoreEventType = 2
	CoreVolume          CoreEventType = 3
	CoreBuffering       CoreEventType = 4
	CorePlaybackUpdate  CoreEventType = 5
	CorePlaybackStart   CoreEventType = 6
	CorePlaybackStop    CoreEventType = 7
	CorePlaybackNearEnd CoreEventType = 8
	CorePlaybackRestart CoreEventType = 9
)

type CoreEvent struct {
	Type CoreEventType
}

type PlayFile struct {
	URI        string
	Start      int32
	Resolution string
	Subtitle   string
	NoAutoPlay bool
	AudioOnly  bool
	VideoOnly  bool
}
