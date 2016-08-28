package httpapi

import (
	"champ/player"
	"champ/plex/model"

	log "github.com/Sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// APIServer contains API server imepmentation
type APIServer struct {
	httpPort        string
	info            *model.Player
	timeline        *Timeline
	playlist        *Playlist
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
	cTimeline := model.NewTimeline("none")
	registy := clientRegistry{items: make(map[string]Subscriber)}
	subscribed := &Subscriptions{&cTimeline, info, &registy}
	playlist := &Playlist{timeline: &cTimeline, info: info, player: engine, subscribed: subscribed}
	go subscribed.Poll()
	go playlist.Poll()

	return &APIServer{
		port, info,
		&Timeline{engine, subscribed, &cTimeline, info},
		playlist,
		&Navigation{},
		engine,
		subscribed,
		&cTimeline,
	}
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

	r.GET("/player/playback/playMedia", a.playlist.PlayMedia)
	r.GET("/player/playback/play", a.playlist.Play)
	r.GET("/player/playback/pause", a.playlist.Pause)
	r.GET("/player/playback/stop", a.playlist.Stop)
	r.GET("/player/playback/setParameters", a.playlist.SetParameters)
	r.GET("/player/playback/seekTo", a.playlist.SeekTo)
	r.GET("/player/playback/stepForward", a.playlist.StepForward)
	r.GET("/player/playback/stepBack", a.playlist.StepBackward)
	r.GET("/player/playback/skipPrevious", a.playlist.SkipPrevious)
	r.GET("/player/playback/skipNext", a.playlist.SkipNext)
	r.GET("/player/navigation/:key", a.navigation.Handle)

	r.Run(":" + a.httpPort)
}
