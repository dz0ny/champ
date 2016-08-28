package httpapi

import (
	"champ/player"
	"champ/plex"
	"champ/plex/model"
	"fmt"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type Playlist struct {
	timeline       *model.Timeline
	info           *model.Player
	player         *player.MPV
	subscribed     *Subscriptions
	current        *model.MediaContainer
	remote         *plex.Remote
	lastUpdateTime time.Time
	ignoreNextStop bool
}

func (a *Playlist) PlayMedia(c *gin.Context) {
	a.player.Stop()
	a.timeline.Clear()
	a.timeline.Type = "none"
	a.subscribed.NotifyAll()
	if containerKey, ok := c.GetQuery("containerKey"); ok {
		machineIdentifier := c.DefaultQuery("machineIdentifier", a.info.MachineIdentifier)
		token := c.DefaultQuery("token", "")
		key := c.DefaultQuery("key", "")
		mediaType := c.DefaultQuery("type", "")
		address := c.DefaultQuery("address", c.ClientIP())
		port := c.DefaultQuery("port", "32400")
		scheme := c.DefaultQuery("protocol", "http")
		a.info.LastCommandID = c.DefaultQuery("commandID", "0")
		lc := log.WithFields(log.Fields{
			"containerKey": containerKey,
			"scheme":       scheme,
			"address":      address,
			"port":         port,
		})
		a.remote = plex.NewRemote(scheme, address, port)
		lc.Info("New Open Request")
		offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

		remoteClient := plex.NewPlexClient(a.remote, nil)
		err, media := remoteClient.GetMedia(containerKey)
		if err == nil {
			a.current = &media
			a.timeline.Protocol = scheme
			a.timeline.Address = address
			a.timeline.Port = port
			a.timeline.MachineIdentifier = machineIdentifier
			a.timeline.Key = key
			a.timeline.ContainerKey = containerKey
			a.timeline.PlayQueueID = media.PlayQueueID
			a.timeline.MediaIndex = media.PlayQueueSelectedItemOffset
			a.timeline.PlayQueueTotalCount = media.PlayQueueTotalCount
			a.timeline.Type = mediaType
			a.timeline.Volume = "100"
			if token != "" {
				a.timeline.Token = token
			}
			a.decideNext("replace", int32(offset/1000))
		} else {
			lc.Error(err)
		}
	}
	c.XML(200, model.Response{Code: "200", Status: "OK"})

}

func (a *Playlist) Poll() {
	a.lastUpdateTime = time.Now()
	for {
		select {
		case event, ok := <-a.player.CoreEventChan:
			if !ok {
				// player has quit, and closed channel
				close(a.player.CoreEventChan)
				return
			}
			switch event.Type {
			case player.CorePlaybackUpdate:
				a.timeline.Time = fmt.Sprintf("%d", a.player.CurrentState.Position)
				if time.Now().After(a.lastUpdateTime.Add(1 * time.Second)) {
					a.lastUpdateTime = time.Now()
					a.subscribed.NotifyAll()
				}
				break
			case player.CoreVolume:
				a.timeline.Volume = fmt.Sprintf("%d", a.player.CurrentState.Volume)
				a.subscribed.NotifyAll()
				break
			case player.CorePause:
				if a.timeline.State == "playing" {
					a.timeline.State = "paused"
					a.subscribed.NotifyAll()
				}
				break
			case player.CoreUnPause:
				if a.timeline.State == "paused" {
					a.timeline.State = "playing"
					a.subscribed.NotifyAll()
				}
				break
			case player.CorePlaybackStart:
				a.timeline.State = "playing"
				a.subscribed.NotifyAll()
				break
			case player.CorePlaybackStop:
				// if next or prev command was issued
				if a.ignoreNextStop {
					a.ignoreNextStop = false
					break
				}
				if a.timeline.MediaIndex+1 == a.timeline.PlayQueueTotalCount {
					a.timeline.Clear()
					a.timeline.Type = "none"
					a.subscribed.NotifyAll()
				} else {
					a.timeline.MediaIndex++
					a.decideNext("replace", 0)
				}
				break
			}

		}
	}

}
func (a *Playlist) decideNext(resolution string, offset int32) {

	a.timeline.Time = ""
	a.timeline.State = "stopped"
	switch a.timeline.Type {
	case "video":
		selected := a.current.Video[a.timeline.MediaIndex]
		a.timeline.PlayQueueItemID = selected.PlayQueueItemID
		a.timeline.Duration = selected.Duration
		a.timeline.SeekRange = fmt.Sprintf("0-%s", a.timeline.Duration)
		a.timeline.RatingKey = selected.RatingKey
		a.timeline.Key = selected.Key
		a.timeline.Guid = selected.Guid
		subURL := ""
		if id, url := selected.SubtitleStream(); url != "" {
			a.timeline.SubtitleStreamID = id
			subURL = a.remote.Path(url)
			log.Infof("Subtitle: %s", subURL)
		}

		mURL := a.remote.Path(selected.VideoStream())
		log.Infof("Video: %s", mURL)
		a.player.Open(&player.PlayFile{URI: mURL, Subtitle: subURL, Resolution: resolution, Start: int32(offset)})
		break
	case "music":
		selected := a.current.Audio[a.timeline.MediaIndex]
		a.timeline.PlayQueueItemID = selected.PlayQueueItemID
		a.timeline.Duration = selected.Duration
		a.timeline.SeekRange = fmt.Sprintf("0-%s", a.timeline.Duration)
		a.timeline.RatingKey = selected.RatingKey
		a.timeline.Guid = selected.Guid

		_, mURL := selected.AudioStream()
		mURL = a.remote.Path(mURL)
		log.Infof("Audio: %s", mURL)
		a.player.Open(&player.PlayFile{URI: mURL, Resolution: resolution, Start: int32(offset), AudioOnly: true})
	}
}

func (a *Playlist) SkipNext(c *gin.Context) {
	a.info.LastCommandID = c.DefaultQuery("commandID", "0")
	a.timeline.MediaIndex++
	if a.timeline.MediaIndex == a.timeline.PlayQueueTotalCount-1 {
		a.timeline.MediaIndex = 0
	}
	a.ignoreNextStop = true
	a.decideNext("replace", 0)
	c.XML(200, model.Response{Code: "200", Status: "OK"})
}

func (a *Playlist) SkipPrevious(c *gin.Context) {
	a.info.LastCommandID = c.DefaultQuery("commandID", "0")
	a.timeline.MediaIndex--
	if a.timeline.MediaIndex <= 0 {
		return
	}
	a.ignoreNextStop = true
	a.decideNext("replace", 0)
	c.XML(200, model.Response{Code: "200", Status: "OK"})
}
func (a *Playlist) Play(c *gin.Context) {
	a.info.LastCommandID = c.DefaultQuery("commandID", "0")
	a.player.Play()
	c.XML(200, model.Response{Code: "200", Status: "OK"})
}

func (a *Playlist) Pause(c *gin.Context) {
	a.info.LastCommandID = c.DefaultQuery("commandID", "0")
	a.player.Pause()
	c.XML(200, model.Response{Code: "200", Status: "OK"})
}

func (a *Playlist) Stop(c *gin.Context) {
	a.info.LastCommandID = c.DefaultQuery("commandID", "0")
	a.player.Stop()
	c.XML(200, model.Response{Code: "200", Status: "OK"})
}

func (a *Playlist) SeekTo(c *gin.Context) {
	a.info.LastCommandID = c.DefaultQuery("commandID", "0")
	log.Info(c.Request.URL.Query())
	if val, ok := c.GetQuery("offset"); ok {
		i, _ := strconv.ParseInt(val, 10, 64)
		a.player.SeekTo(int32(i / 1000))
	}
	c.XML(200, model.Response{Code: "200", Status: "OK"})
}

func (a *Playlist) StepForward(c *gin.Context) {
	a.info.LastCommandID = c.DefaultQuery("commandID", "0")
	log.Info(c.Request.URL.Query())
	a.player.SeekTo(int32(a.player.CurrentState.Position/1000) + 30)
	c.XML(200, model.Response{Code: "200", Status: "OK"})
}

func (a *Playlist) StepBackward(c *gin.Context) {
	a.info.LastCommandID = c.DefaultQuery("commandID", "0")
	log.Info(c.Request.URL.Query())
	a.player.SeekTo(int32(a.player.CurrentState.Position/1000) - 15)
	c.XML(200, model.Response{Code: "200", Status: "OK"})
}

func (a *Playlist) SetParameters(c *gin.Context) {
	a.info.LastCommandID = c.DefaultQuery("commandID", "0")
	log.Info(c.Request.URL.Query())
	if val, ok := c.GetQuery("volume"); ok {
		i, _ := strconv.ParseInt(val, 10, 64)
		a.player.SetVolume(int32(i))
	}
	c.XML(200, model.Response{Code: "200", Status: "OK"})
}
