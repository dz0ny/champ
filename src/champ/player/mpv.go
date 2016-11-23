package player

/*
#include <mpv/client.h>
#include <stdlib.h>
#cgo LDFLAGS: -lmpv
*/
import "C"
import (
	"fmt"
	"strings"
	"unsafe"

	log "github.com/Sirupsen/logrus"

	"github.com/YouROK/go-mpv/mpv"
)

type MPV struct {
	backend       *mpv.Mpv
	eventChan     chan *mpv.Event
	CurrentState  *PlayState
	CoreEventChan chan *CoreEvent
}

func (m *MPV) Initialize() error {
	m.backend = mpv.Create()
	m.eventChan = make(chan *mpv.Event)
	m.CoreEventChan = make(chan *CoreEvent)
	m.CurrentState = &PlayState{}
	go func() {
		for {
			e := m.backend.WaitEvent(-1)
			m.eventChan <- e
		}
	}()
	log.Debugln("Set options")
	platformInitialize(m.backend)
	m.backend.ObserveProperty(0, "seeking", mpv.FORMAT_FLAG)
	m.backend.ObserveProperty(0, "core-idle", mpv.FORMAT_FLAG)
	m.backend.ObserveProperty(0, "cache-buffering-state", mpv.FORMAT_INT64)
	m.backend.ObserveProperty(0, "playback-time", mpv.FORMAT_DOUBLE)
	m.backend.ObserveProperty(0, "duration", mpv.FORMAT_DOUBLE)
	m.backend.ObserveProperty(0, "volume", mpv.FORMAT_DOUBLE)
	log.Debugln("Initialize")
	return m.backend.Initialize()
}

func (m *MPV) Loop() {
	log.Debugln("Start loop")
	for {
		select {
		case event, ok := <-m.eventChan:
			if !ok {
				// player has quit, and closed channel
				close(m.eventChan)
				return
			}
			m.processMPVEvent(event)
		}
	}
}

func (m *MPV) Open(f *PlayFile) {
	if f.URI == "" {
		return
	}

	cmd := []string{
		"loadfile", f.URI,
	}
	if f.Resolution != "" {
		cmd = append(cmd, f.Resolution)
	}
	ecmd := []string{}
	if f.Subtitle != "" {
		ecmd = append(ecmd, fmt.Sprintf("sub-file=%s", f.Subtitle), "sid=auto")
	}
	if f.Start > 0 {
		ecmd = append(ecmd, fmt.Sprintf("start=+%d", f.Start))
	}
	if f.NoAutoPlay {
		ecmd = append(ecmd, "pause=yes")
	}
	if f.VideoOnly && !f.AudioOnly {
		ecmd = append(ecmd, "aid=no")
	}
	if !f.VideoOnly && f.AudioOnly {
		ecmd = append(ecmd, "vid=no")
	}
	ecmd = append(ecmd, "ad=''", "vd=''")
	cmd = append(cmd, strings.Join(ecmd, ","))
	log.Debugln(cmd)
	m.backend.Command(cmd)
}

func (m *MPV) Pause() {
	switch m.CurrentState.State {
	case STATE_PLAYING:
		m.backend.Command([]string{"set", "pause", "yes"})
	}
}

func (m *MPV) Play() {
	switch m.CurrentState.State {
	case STATE_STOPPED, STATE_PAUSED:
		m.backend.Command([]string{"set", "pause", "no"})
	}
}

func (m *MPV) SetVolume(volume int32) {
	if volume > 100 {
		return
	}
	m.backend.Command([]string{"set", "volume", fmt.Sprintf("%d", volume)})
}

func (m *MPV) SeekTo(location int32) {
	switch m.CurrentState.State {
	case STATE_PLAYING:
		m.backend.Command([]string{"seek", fmt.Sprintf("%d", location), "absolute+exact"})
	case STATE_PAUSED:
		m.backend.Command([]string{"seek", fmt.Sprintf("%d", location), "absolute+exact"})
		m.backend.Command([]string{"set", "pause", "no"})
	}

}

func (m *MPV) Stop() {
	m.backend.Command([]string{"stop"})
}

func (m *MPV) processMPVPropertyChange(e *mpv.Event) {
	prop := (*C.mpv_event_property)(e.Data.(unsafe.Pointer))
	if mpv.Format(prop.format) == mpv.FORMAT_NONE {
		return
	}
	propName := C.GoString(prop.name)
	switch propName {
	case "seeking":
		if *(*bool)(prop.data) {
			m.CurrentState.State = STATE_SEEKING
		}
		m.CoreEventChan <- &CoreEvent{CoreSeek}
		break
	case "volume":
		m.CurrentState.Volume = int32(*(*float64)(prop.data))
		m.CoreEventChan <- &CoreEvent{CoreVolume}
		break
	case "cache-buffering-state":
		fillPercent := *(*int64)(prop.data)
		if fillPercent < 100 {
			m.CurrentState.State = STATE_BUFFERING
			m.CoreEventChan <- &CoreEvent{CoreBuffering}
		}

		break
	case "playback-time":
		pos := int32(*(*float64)(prop.data) * 1000)
		if pos != m.CurrentState.Position {
			m.CurrentState.Position = pos
			m.CurrentState.State = STATE_PLAYING
			m.CoreEventChan <- &CoreEvent{CorePlaybackUpdate}
			if pos+5000 == m.CurrentState.Duration { //5s
				m.CoreEventChan <- &CoreEvent{CorePlaybackNearEnd}
			}
		}
		break
	case "duration":
		m.CurrentState.Duration = int32(*(*float64)(prop.data) * 1000)
		m.CoreEventChan <- &CoreEvent{CorePlaybackUpdate}
		break
	case "core-idle":
		if *(*bool)(prop.data) {
			m.CurrentState.State = STATE_STOPPED
			m.CurrentState.Duration = 0
			m.CurrentState.Position = 0
			m.CoreEventChan <- &CoreEvent{CoreIdle}
		}
	}

}

func (m *MPV) processMPVEvent(e *mpv.Event) {
	if e.Event_Id != mpv.EVENT_NONE && e.Event_Id != mpv.EVENT_LOG_MESSAGE && e.Event_Id != mpv.EVENT_PROPERTY_CHANGE {
		log.Infoln(e.Event_Id.String())
	}
	switch e.Event_Id {
	case mpv.EVENT_IDLE:
		m.CurrentState.State = STATE_STOPPED
		m.CurrentState.Duration = 0
		m.CurrentState.Position = 0
		m.CoreEventChan <- &CoreEvent{CoreIdle}
	case mpv.EVENT_PLAYBACK_RESTART:
		m.CoreEventChan <- &CoreEvent{CorePlaybackRestart}
		break
	case mpv.EVENT_NONE, mpv.EVENT_LOG_MESSAGE:
		break
	case mpv.EVENT_END_FILE:
		m.CurrentState.State = STATE_STOPPED
		m.CurrentState.Duration = 0
		m.CurrentState.Position = 0
		m.CoreEventChan <- &CoreEvent{CorePlaybackStop}
		break
	case mpv.EVENT_START_FILE:
		m.CurrentState.State = STATE_PLAYING
		m.CurrentState.Duration = 0
		m.CurrentState.Position = 0
		m.CoreEventChan <- &CoreEvent{CorePlaybackStart}
		break
	case mpv.EVENT_PAUSE:
		m.CurrentState.State = STATE_PAUSED
		m.CoreEventChan <- &CoreEvent{CorePause}
		break
	case mpv.EVENT_UNPAUSE:
		m.CurrentState.State = STATE_PLAYING
		m.CoreEventChan <- &CoreEvent{CoreUnPause}
		break
	case mpv.EVENT_PROPERTY_CHANGE:
		m.processMPVPropertyChange(e)
		break
	default:
		log.Debugln("processMPVEvent", e.Event_Id.String())
	}

}
