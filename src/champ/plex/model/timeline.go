package model

// <Timeline
// state="playing"
// time="122530"
// type="video"

// location="fullScreenVideo"
// key="/library/metadata/3698"
// ratingKey="3698"
// containerKey="/library/metadata/3698"

// duration="3319906"
// seekRange="0-3319906"
// controllable="volume,shuffle,repeat,audioStream,videoStream,subtitleStream,skipPrevious,skipNext,seekTo,stepBack,stepForward,stop,playPause"
// machineIdentifier="dd4c35b96c82136b6660ca39ebf1d1843b53e24d"

// protocol="https"
// address="192.168.2.33"
// port="32400"
// guid="com.plexapp.agents.none://39c50223ad12e1b343122ee54ba1d6730f364333?lang=xn"

// volume="100"
// shuffle="0"
// mute="0"
// repeat="0"

// subtitleStreamID="-1"
// audioStreamID="-1"
// />

// <Timeline state="stopped" time="0" type="music" seekRange="0-0" />
type Timeline struct {
	State            string `xml:"state,attr"`
	Time             string `xml:"time,attr,omitempty"`
	Type             string `xml:"type,attr"`
	PlayQueueVersion string `xml:"playQueueVersion,attr"`

	Location     string `xml:"location,attr,omitempty"`
	Key          string `xml:"key,attr,omitempty"`
	RatingKey    string `xml:"ratingKey,attr,omitempty"`
	ContainerKey string `xml:"containerKey,attr,omitempty"`

	PlayQueueID         int `xml:"playQueueID,attr,omitempty"`
	PlayQueueItemID     int `xml:"playQueueItemID,attr,omitempty"`
	MediaIndex          int `xml:"mediaIndex,attr,omitempty"`
	PlayQueueTotalCount int `xml:"playQueueTotalCount,attr,omitempty"`

	Duration          string `xml:"duration,attr,omitempty"`
	SeekRange         string `xml:"seekRange,attr,omitempty"`
	Controllable      string `xml:"controllable,attr,omitempty"`
	MachineIdentifier string `xml:"machineIdentifier,attr,omitempty"`
	Token             string `xml:"token,attr,omitempty"`

	Protocol string `xml:"protocol,attr,omitempty"`
	Address  string `xml:"address,attr,omitempty"`
	Port     string `xml:"port,attr,omitempty"`
	Guid     string `xml:"guid,attr,omitempty"`

	Volume string `xml:"volume,attr,omitempty"`

	SubtitleStreamID string `xml:"subtitleStreamID,attr,omitempty"`
	AudioStreamID    string `xml:"audioStreamID,attr,omitempty"`
}

func NewTimeline(timelineType string) Timeline {
	con := ""
	switch timelineType {
	case "photo":
		con = "skipPrevious,skipNext,stop"
		break
	case "music":
		con = "playPause,stop,volume,shuffle,repeat,seekTo,skipPrevious,skipNext,stepBack,stepForward"
		break
	case "video":
		con = "playPause,stop,volume,audioStream,subtitleStream,seekTo,skipPrevious,skipNext,stepBack,stepForward"
		break
	}
	return Timeline{
		State:            "stopped",
		PlayQueueVersion: "1",
		Type:             timelineType,
		Controllable:     con,
	}
}

func (t *Timeline) Clear() {
	t.State = "stopped"
	t.Time = ""
	t.PlayQueueVersion = "1"

	t.Location = ""
	t.Key = ""
	t.RatingKey = ""
	t.ContainerKey = ""

	t.PlayQueueID = -1
	t.PlayQueueItemID = -1
	t.MediaIndex = -1
	t.PlayQueueTotalCount = -1

	t.Duration = ""
	t.SeekRange = ""
	t.MachineIdentifier = ""

	t.Protocol = ""
	t.Address = ""
	t.Port = ""
	t.Guid = ""

	t.Volume = ""

	t.SubtitleStreamID = ""
	t.AudioStreamID = ""
}
