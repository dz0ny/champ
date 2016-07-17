package httpapi

import (
	"champ/player"
	"champ/plex"
	"champ/plex/model"
	"fmt"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

type Playback struct {
	player          *player.MPV
	currentTimeline *model.Timeline
	info            *model.Player
	subs            *Subscriptions
}

func getPath(p string) string {
	return strings.Split(p, "?")[0]
}

func getPlayQueueID(p string) string {
	l := strings.Split(p, "/")
	return l[len(l)-1]
}

func (a *Playback) PlayMedia(c *gin.Context) {
	log.Debug(c.Request.URL.Query())
	//map[port:[32400] containerKey:[/playQueues/967?own=1&repeat=0&window=200] type:[video] protocol:[http] address:[192.168.2.33] key:[/library/metadata/3698] offset:[0] commandID:[2] machineIdentifier:[dd4c35b96c82136b6660ca39ebf1d1843b53e24d]]
	if key, ok := c.GetQuery("key"); ok {
		machineIdentifier := c.DefaultQuery("machineIdentifier", a.info.MachineIdentifier)
		containerKey := c.DefaultQuery("containerKey", "")
		token := c.DefaultQuery("token", "")
		address := c.DefaultQuery("address", c.ClientIP())
		port := c.DefaultQuery("port", "32400")
		scheme := c.DefaultQuery("protocol", "http")
		a.info.LastCommandID = c.DefaultQuery("commandID", "0")
		lc := log.WithFields(log.Fields{
			"scheme":  scheme,
			"address": address,
			"port":    port,
		})
		remote := plex.NewRemote(scheme, address, port)
		lc.Info("New Open Request")
		offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

		remoteClient := plex.NewPlexClient(remote, nil)
		err, media := remoteClient.GetMedia(key)
		if err == nil {
			// a.currentTimeline.SubtitleStreamID = "-1"
			// a.currentTimeline.AudioStreamID = "-1"
			if containerKey != "" {
				qp := getPath(containerKey)
				a.currentTimeline.ContainerKey = qp
				a.currentTimeline.PlayQueueID = getPlayQueueID(qp)
			}
			if token != "" {
				a.currentTimeline.Token = token
			}

			a.currentTimeline.Protocol = scheme
			a.currentTimeline.Address = address
			a.currentTimeline.Port = port
			a.currentTimeline.MachineIdentifier = machineIdentifier
			a.currentTimeline.Key = key

			if media.Video != nil {
				a.currentTimeline.Type = "video"
				a.currentTimeline.Duration = media.Video.Duration
				a.currentTimeline.SeekRange = fmt.Sprintf("0-%s", media.Video.Duration)
				a.currentTimeline.RatingKey = media.Video.RatingKey
				a.currentTimeline.Guid = media.Video.Guid

				subURL := ""
				if media.SubtitleStream() != "" {
					subURL = remote.Path(media.SubtitleStream())
					lc.Infof("Subtitle: %s", subURL)
				}
				mURL := remote.Path(media.VideoStream())
				lc.Infof("Video: %s", mURL)
				a.player.Open(&player.PlayFile{URI: mURL, Subtitle: subURL, Resolution: "replace", Start: int32(offset / 1000)})
			}
			if media.Audio != nil {
				a.currentTimeline.Type = "music"
				a.currentTimeline.Duration = media.Audio.Duration
				a.currentTimeline.SeekRange = fmt.Sprintf("0-%s", media.Audio.Duration)
				a.currentTimeline.RatingKey = media.Audio.RatingKey
				a.currentTimeline.Guid = media.Audio.Guid
				mURL := remote.Path(media.AudioStream())
				lc.Infof("Audio: %s", mURL)
				a.player.Open(&player.PlayFile{URI: mURL, Resolution: "replace", Start: int32(offset / 1000), AudioOnly: true})
			}

		} else {
			lc.Error(err)
		}
	}
	c.XML(200, model.Response{Code: "200", Status: "OK"})

}

func (a *Playback) Play(c *gin.Context) {
	a.info.LastCommandID = c.DefaultQuery("commandID", "0")
	a.player.Play()
	c.XML(200, model.Response{Code: "200", Status: "OK"})
}

func (a *Playback) Pause(c *gin.Context) {
	a.info.LastCommandID = c.DefaultQuery("commandID", "0")
	a.player.Pause()
	c.XML(200, model.Response{Code: "200", Status: "OK"})
}

func (a *Playback) Stop(c *gin.Context) {
	a.info.LastCommandID = c.DefaultQuery("commandID", "0")
	a.player.Stop()
	c.XML(200, model.Response{Code: "200", Status: "OK"})
}

func (a *Playback) SeekTo(c *gin.Context) {
	a.info.LastCommandID = c.DefaultQuery("commandID", "0")
	log.Info(c.Request.URL.Query())
	if val, ok := c.GetQuery("offset"); ok {
		i, _ := strconv.ParseInt(val, 10, 64)
		a.player.SeekTo(int32(i / 1000))
	}
	c.XML(200, model.Response{Code: "200", Status: "OK"})
}

func (a *Playback) StepForward(c *gin.Context) {
	a.info.LastCommandID = c.DefaultQuery("commandID", "0")
	log.Info(c.Request.URL.Query())
	a.player.SeekTo(int32(a.player.CurrentState.Position/1000) + 30)
	c.XML(200, model.Response{Code: "200", Status: "OK"})
}

func (a *Playback) StepBackward(c *gin.Context) {
	a.info.LastCommandID = c.DefaultQuery("commandID", "0")
	log.Info(c.Request.URL.Query())
	a.player.SeekTo(int32(a.player.CurrentState.Position/1000) - 15)
	c.XML(200, model.Response{Code: "200", Status: "OK"})
}

func (a *Playback) SetParameters(c *gin.Context) {
	a.info.LastCommandID = c.DefaultQuery("commandID", "0")
	log.Info(c.Request.URL.Query())
	if val, ok := c.GetQuery("volume"); ok {
		i, _ := strconv.ParseInt(val, 10, 64)
		a.player.SetVolume(int32(i))
	}
	c.XML(200, model.Response{Code: "200", Status: "OK"})
}
