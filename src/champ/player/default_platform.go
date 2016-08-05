// +build !rpi

package player

import "github.com/YouROK/go-mpv/mpv"

func platformInitialize(m *mpv.Mpv) {
	m.RequestLogMessages("terminal-default")
	m.SetOptionString("msg-level", "all=v")
	m.SetOptionString("osd-level", "0")
	m.SetOptionString("input-cursor", "no")
	m.SetOptionString("cursor-autohide", "no")
	m.SetOptionString("hwdec-preload", "auto")
	m.SetOptionString("softvol", "yes")
	m.SetOptionString("gapless-audio", "yes")
	// https://ffmpeg.org/ffmpeg-filters.html#replaygain
	m.SetOptionString("af", "volume=replaygain=album")
	m.SetOptionString("audio-client-name", "ChampPlayer")
	m.SetOptionString("title", "${?media-title:${media-title}}${!media-title:No file.}")

	m.SetOptionString("config", "yes")
	m.SetOptionString("config-dir", "/etc/champ")

	// Audio defaults
	// https://github.com/mpv-player/mpv/blob/master/DOCS/man/af.rst
	m.SetOptionString("af-defaults", "lavrresample:o=[surround_mix_level=1]:normalize=yes")

	m.SetOption("ad-lavc-downmix", mpv.FORMAT_FLAG, false)
	m.SetOption("demuxer-mkv-probe-start-time", mpv.FORMAT_FLAG, false)
	m.SetOption("no-resume-playback", mpv.FORMAT_FLAG, true)
	m.SetOption("no-input-terminal", mpv.FORMAT_FLAG, true)
	m.SetOption("fullscreen", mpv.FORMAT_FLAG, true)
	m.SetOption("input-media-keys", mpv.FORMAT_FLAG, true)

	//cache
	m.SetOption("cache-default", mpv.FORMAT_INT64, 160) // 10 seconds
	m.SetOption("cache-seek-min", mpv.FORMAT_INT64, 16) // 1 second
	m.SetOption("cache-seek-min", mpv.FORMAT_INT64, 16) // 1 second

}
