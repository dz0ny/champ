package main

import (
	"champ/player"
	"os"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var build = "dev"

var app = kingpin.New("spinwheel", "Shuffle player which also plays from YouTube(tm)")
var showVerbose = app.Flag("debug", "Verbose mode.").Bool()
var configFile = app.Flag("playlist", "Path to playlist file (.yaml)").Short('p').ExistingFile()

func main() {
	app.Author("dz0ny")
	app.Version(build)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *showVerbose {
		log.SetLevel(log.DebugLevel)
		log.Info("Debug mode enabled")
	}
	videos := LoadVideos(*configFile)
	p := &player.MPV{}

	p.Initialize()
	// Start the mainloop.
	go p.Loop()
	for {
		select {
		case event, ok := <-p.CoreEventChan:
			if !ok {
				// player has quit, and closed channel
				close(p.CoreEventChan)
				return
			}
			switch event.Type {
			case player.CoreIdle, player.CorePlaybackNearEnd:
				p.Open(&player.PlayFile{URI: RandomVideo(videos), Resolution: "append-play"})
				break
			}

		}
	}
}
