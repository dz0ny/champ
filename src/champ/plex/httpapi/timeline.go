package httpapi

import (
	"champ/player"
	"champ/plex"
	"champ/plex/model"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

type clientRegistry struct {
	sync.Mutex
	items map[string]Subscriber
}

type Timeline struct {
	player          *player.MPV
	clients         *Subscriptions
	currentTimeline *model.Timeline
	info            *model.Player
}

type Subscriber struct {
	name          string // X-Plex-Device-Name
	lastSubscribe time.Time
	remote        *plex.Remote
}

type Subscriptions struct {
	currentTimeline *model.Timeline
	info            *model.Player
	clientRegistry  *clientRegistry
}

func getTimelines(currentTimeline *model.Timeline) []model.Timeline {
	t := []model.Timeline{
		model.NewTimeline("photo"),
	}
	switch currentTimeline.Type {
	case "video":
		t = append(t, *currentTimeline, model.NewTimeline("music"))
		break
	case "music":
		t = append(t, *currentTimeline, model.NewTimeline("video"))
		break
	default:
		t = append(t, model.NewTimeline("music"), model.NewTimeline("video"))
	}
	return t
}

func (s *Subscriptions) Add(key, name string, rem *plex.Remote) {
	s.clientRegistry.Lock()
	alog := log.WithFields(log.Fields{
		"key":    key,
		"name":   name,
		"remote": rem,
	})
	if r, found := s.clientRegistry.items[key]; found {
		r.lastSubscribe = time.Now()
		alog.Debug("Updating subscription")
	} else {
		s.clientRegistry.items[key] = Subscriber{name, time.Now(), rem}
		alog.Debug("New subscription")
	}
	s.clientRegistry.Unlock()
}

func (s *Subscriptions) Remove(key string) {
	s.clientRegistry.Lock()
	if _, found := s.clientRegistry.items[key]; found {
		delete(s.clientRegistry.items, key)
	}
	s.clientRegistry.Unlock()
}

func (s *Subscriptions) Poll() {

	//Player Subscription times out after 90 seconds.
	s.clientRegistry.Lock()
	for name, client := range s.clientRegistry.items {
		if client.lastSubscribe.After(time.Now().Add(-90 * time.Second)) {
			log.WithFields(log.Fields{
				"name":          name,
				"lastSubscribe": client.lastSubscribe,
			}).Debug("Removing from subscription registry")
			s.Remove(name)
		}
	}
	s.clientRegistry.Unlock()
}

func (s *Subscriptions) Notify(r *plex.Remote) {
	data := model.MediaContainer{
		MachineIdentifier: s.info.MachineIdentifier,
		Timelines:         getTimelines(s.currentTimeline),
		CommandID:         s.info.LastCommandID,
	}
	r.Notify(data)
}

func (s *Subscriptions) NotifyAll() {
	s.clientRegistry.Lock()
	for _, client := range s.clientRegistry.items {
		s.Notify(client.remote)
	}
	s.clientRegistry.Unlock()
}

func (a *Timeline) Subscribe(c *gin.Context) {
	// Players MUST accept and process commands received from unsubscribed controllers, even if they lack a commandID.
	if commandID := c.DefaultQuery("commandID", ""); commandID != "" {
		a.info.LastCommandID = commandID
	}
	address := c.ClientIP()
	port := c.DefaultQuery("port", "32400")
	scheme := c.DefaultQuery("protocol", "http")
	ci := c.Request.Header.Get("X-Plex-Client-Identifier")
	cn := c.Request.Header.Get("X-Plex-Device-Name")
	if ci != "" && cn != "" {
		r := plex.NewRemote(scheme, address, port)
		a.clients.Notify(r)
		a.clients.Add(ci, cn, r)
		c.XML(200, model.Response{Code: "200", Status: "OK"})
	} else {
		c.XML(404, model.Response{Code: "404", Status: "Not Found - Plex Client"})
	}
}

func (a *Timeline) UnSubscribe(c *gin.Context) {
	ci := c.Request.Header.Get("X-Plex-Client-Identifier")
	if ci != "" {
		a.clients.Remove(ci)
	}
	c.XML(200, model.Response{Code: "200", Status: "OK"})
}

func (a *Timeline) Poll(c *gin.Context) {
	commandID := c.DefaultQuery("commandID", a.info.LastCommandID)

	if wait := c.DefaultQuery("wait", "0"); wait == "1" {
		time.Sleep(30 * time.Second)
	}
	t := getTimelines(a.currentTimeline)
	c.XML(200, model.MediaContainer{MachineIdentifier: a.info.MachineIdentifier, Timelines: t, CommandID: commandID})
}
