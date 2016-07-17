package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/parnurzeal/gorequest"

	"github.com/otium/ytdl"
)

const atvURL = "http://a1.phobos.apple.com/us/r1000/000/Features/atv/AutumnResources/videos/entries.json"

func init() {
	rand.Seed(time.Now().UnixNano())
}

type ATV []struct {
	Assets []struct {
		URL         string `json:"url"`
		Description string `json:"accessibilityLabel"`
		TimeOfDay   string `json:"timeOfDay"`
	} `json:"assets"`
}

func RandomATV() string {
	var a = ATV{}
	gorequest.New().Get(atvURL).EndStruct(&a)
	l := []string{}
	for _, asset := range a {
		for _, s := range asset.Assets {
			l = append(l, s.URL)
			fmt.Printf("\"%s\",  // %s - %s\n", s.URL, s.Description, s.TimeOfDay)
		}
	}
	return l[rand.Intn(len(l))]
}

type Playlist struct {
	Videos []string `yaml:"videos"`
}

func LoadVideos(path string) []string {
	var t Playlist
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &t)
	if err != nil {
		panic(err)
	}

	return t.Videos
}

func RandomVideo(videos []string) string {
	video := videos[rand.Intn(len(videos))]
	if strings.Contains(video, "youtube.com") {
		log.Infof("Resolving video %s", video)
		vid, err := ytdl.GetVideoInfo(video)
		if err != nil {
			log.Errorln(err)
			return RandomATV()
		}
		format := vid.Formats.Best(ytdl.FormatItagKey)
		videourl, err := vid.GetDownloadURL(format[0])
		if err != nil {
			log.Errorln(err)
			return RandomATV()
		}
		video = videourl.String()
	}
	return video
}
