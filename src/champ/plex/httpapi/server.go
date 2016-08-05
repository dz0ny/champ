package httpapi

import (
	"champ/player"
	"champ/plex/model"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// APIServer contains API server imepmentation
type APIServer struct {
	httpPort        string
	info            *model.Player
	playback        *Playback
	timeline        *Timeline
	navigation      *Navigation
	engine          *player.MPV
	subs            *Subscriptions
	currentTimeline *model.Timeline
}

// NewAPIServer creates new API server instance
func NewAPIServer(engine *player.MPV, info *model.Player, port string, debug bool) *APIServer {
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	lastUpdateTime := time.Now()
	cTimeline := model.NewTimeline("none")
	registy := clientRegistry{items: make(map[string]Subscriber)}
	subs := &Subscriptions{currentTimeline: &cTimeline, info: info, clientRegistry: &registy}
	go subs.Poll()
	go func(t *model.Timeline, p *player.MPV) {
		for {
			select {
			case event, ok := <-p.CoreEventChan:
				if !ok {
					// player has quit, and closed channel
					close(p.CoreEventChan)
					return
				}
				switch event.Type {
				case player.CorePlaybackUpdate:
					t.Time = fmt.Sprintf("%d", engine.CurrentState.Position)
					if time.Now().After(lastUpdateTime.Add(1 * time.Second)) {
						log.Debugln(engine.CurrentState)
						lastUpdateTime = time.Now()
						subs.NotifyAll()
					}
					break
				case player.CoreVolume:
					t.Volume = fmt.Sprintf("%d", engine.CurrentState.Volume)
					subs.NotifyAll()
					break
				case player.CorePlaybackRestart:
					subs.NotifyAll()
					break
				case player.CorePause:
					if t.State == "playing" {
						t.State = "paused"
						subs.NotifyAll()
					}
					break
				case player.CorePlaybackStart:
					t.State = "playing"
					subs.NotifyAll()
					break
				case player.CorePlaybackStop:
					t.Clear()
					t.Type = "none"
					subs.NotifyAll()
					break
				}

			}
		}
	}(&cTimeline, engine)

	return &APIServer{port, info, &Playback{engine, &cTimeline, info, subs}, &Timeline{engine, subs, &cTimeline, info}, &Navigation{}, engine, subs, &cTimeline}
}

func (a *APIServer) getResources(c *gin.Context) {
	log.Info(c.Request.URL.Query())
	c.XML(200, model.MediaContainer{Player: a.info})
}

func NewPlexCORS(info *model.Player) gin.HandlerFunc {
	return func(c *gin.Context) {
		clog := log.WithField("query", c.Request.URL.Query()).WithField("remote", c.ClientIP())
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Add("Access-Control-Max-Age", "1209600")

		rh := c.Request.Header.Get("access-control-request-headers")
		if rh != "" {
			c.Writer.Header().Add("Access-Control-Allow-Headers", rh)
		}
		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Add("Content-Type", "text/plain")
			c.Writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT, HEAD")

			c.Writer.Header().Add("Connection", "close")
			c.AbortWithStatus(200)
		} else {
			c.Writer.Header().Add("X-Plex-Client-Identifier", info.MachineIdentifier)
			c.Writer.Header().Add("X-Plex-Protocol", "1.0")
			clog.Debug("Command request")
		}

		c.Next()
	}
}

// ListenAndServe exposes internal server
func (a *APIServer) ListenAndServe() {
	r := gin.Default()
	r.Use(NewPlexCORS(a.info))
	r.GET("/resources", a.getResources)

	r.GET("/player/timeline/subscribe", a.timeline.Subscribe)
	r.GET("/player/timeline/unsubscribe", a.timeline.UnSubscribe)
	r.GET("/player/timeline/poll", a.timeline.Poll)

	r.GET("/player/playback/playMedia", a.playback.PlayMedia)
	r.GET("/player/playback/play", a.playback.Play)
	r.GET("/player/playback/pause", a.playback.Pause)
	r.GET("/player/playback/stop", a.playback.Stop)
	r.GET("/player/playback/setParameters", a.playback.SetParameters)
	r.GET("/player/playback/seekTo", a.playback.SeekTo)
	r.GET("/player/playback/stepForward", a.playback.StepForward)
	r.GET("/player/playback/stepBack", a.playback.StepBackward)

	r.GET("/player/navigation/:key", a.navigation.Handle)

	r.Run(":" + a.httpPort)
}
