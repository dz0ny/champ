package main

import (
	"champ/player"
	"champ/plex/gdm"
	"champ/plex/httpapi"
	"champ/plex/model"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"
	"gopkg.in/alecthomas/kingpin.v2"
)

var build = "dev"

var app = kingpin.New("champ", "Minimalistic Plex 2nd screen client")
var showVerbose = app.Flag("debug", "Verbose mode.").Bool()
var title = app.Flag("name", "Name of this player").Short('n').Default("Champ Player").String()
var httpPort = app.Flag("port", "HTTP server port").Short('p').Default("32016").String()

func getUUID() string {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range interfaces {
			if len(iface.HardwareAddr) >= 6 {
				return uuid.NewV3(uuid.NamespaceDNS, string(iface.HardwareAddr)+*title).String()
			}
		}
	}
	return uuid.NewV1().String()
}

func main() {
	app.Author("dz0ny")
	app.Version(build)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *showVerbose {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug mode enabled")
	}
	uid := getUUID()
	log.WithField("UUID", uid).Info("Starting")
	plexPlayer := model.Player{
		Title:                *title,
		Protocol:             "plex",
		ProtocolVersion:      "1",
		ProtocolCapabilities: "navigation,playback,timeline,playqueues",
		MachineIdentifier:    uid,
		Product:              "Champ Player",
		Platform:             runtime.GOOS,
		PlatformVersion:      runtime.Version(),
		DeviceClass:          "sbc",
	}

	engine := &player.MPV{}
	engine.Initialize()
	// Start the mainloop.
	go engine.Loop()

	_, s := gdm.NewAdvertiser(&plexPlayer, *httpPort)
	defer s.Shutdown()

	api := httpapi.NewAPIServer(engine, &plexPlayer, *httpPort, *showVerbose)
	go api.ListenAndServe()

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch //block

}
