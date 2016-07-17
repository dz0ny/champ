// +build rpi

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

	// The backbuffer makes seeking back faster (without having to do a HTTP-level seek)
	m.SetOption("cache-backbuffer", mpv.FORMAT_INT64, 10*1024) // KB
	// The demuxer queue is used for the readahead, and also for dealing with badly
	// interlaved audio/video. Setting it too low increases sensitivity to network
	// issues, and could cause playback failure with "bad" files.
	m.SetOption("demuxer-max-bytes", mpv.FORMAT_INT64, 50*1024*1024) // bytes
	// Specifically for enabling mpeg4.
	m.SetOptionString("hwdec-codecs", "all")
	// Do not use exact seeks by default. (This affects the start position in the "loadfile"
	// command in particular. We override the seek mode for normal "seek" commands.)
	m.SetOptionString("hr-seek", "no")

	// Force RPi video engine
	m.SetOptionString("vo", "rpi:background=yes")
}
