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
	CoreIdle            CoreEventType = 0
	CorePause           CoreEventType = 1
	CoreUnPause         CoreEventType = 2
	CoreSeek            CoreEventType = 3
	CoreVolume          CoreEventType = 4
	CoreBuffering       CoreEventType = 5
	CorePlaybackUpdate  CoreEventType = 6
	CorePlaybackStart   CoreEventType = 7
	CorePlaybackStop    CoreEventType = 8
	CorePlaybackNearEnd CoreEventType = 9
	CorePlaybackRestart CoreEventType = 10
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
