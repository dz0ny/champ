package model

import (
	"encoding/xml"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVideoService(t *testing.T) {

	var xmls = `
  <MediaContainer size="1" allowSync="1" identifier="com.plexapp.plugins.library" librarySectionID="3" librarySectionTitle="Movies" librarySectionUUID="4664608e-e7fc-47c5-a8ce-f87788721ec6" mediaTagPrefix="/system/bundle/media/flags/" mediaTagVersion="1461959475">
<Video ratingKey="2593" key="/library/metadata/2593" guid="com.plexapp.agents.imdb://tt2488496?lang=en" librarySectionID="3" studio="Lucasfilm" type="movie" title="Star Wars: The Force Awakens" contentRating="PG-13" summary="Thirty years after defeating the Galactic Empire, Han Solo and his allies face a new threat from the evil Kylo Ren and his army of Stormtroopers." rating="7.6" viewOffset="7200307" lastViewedAt="1466843725" year="2015" tagline="Every generation has a story." thumb="/library/metadata/2593/thumb/1465237688" art="/library/metadata/2593/art/1465237688" duration="8286932" originallyAvailableAt="2015-12-14" addedAt="1465237566" updatedAt="1465237688" chapterSource="">
<Media videoResolution="720" id="2474" duration="8286932" bitrate="1049" width="1280" height="536" aspectRatio="2.35" audioChannels="2" audioCodec="aac" videoCodec="h264" container="mp4" videoFrameRate="24p" optimizedForStreaming="1" audioProfile="lc" has64bitOffsets="0" videoProfile="high">
<Part id="2478" key="/library/parts/2478/file.mp4" duration="8286932" file="Star Wars The Force Awakens.mp4" size="1086490605" audioProfile="lc" container="mp4" has64bitOffsets="0" optimizedForStreaming="1" videoProfile="high">
<Stream id="5887" streamType="1" default="1" codec="h264" index="0" bitrate="949" bitDepth="8" cabac="1" chromaSubsampling="4:2:0" codecID="avc1" colorRange="tv" colorSpace="bt709" duration="8286695" frameRate="23.976" frameRateMode="cfr" hasScalingMatrix="0" height="536" level="41" pixelFormat="yuv420p" profile="high" refFrames="5" scanType="progressive" streamIdentifier="1" width="1280" />
<Stream id="5888" streamType="2" selected="1" default="1" codec="aac" index="1" channels="2" bitrate="96" audioChannelLayout="stereo" bitrateMode="vbr" codecID="40" duration="8286932" profile="lc" samplingRate="48000" streamIdentifier="2" />
<Stream id="5889" key="/library/streams/5889" streamType="3" codec="srt" format="srt" />
</Part>
</Media>
</Video>
</MediaContainer>

`

	Convey("Test Getting elements", t, func() {
		m := MediaContainer{}
		xml.Unmarshal([]byte(xmls), &m)
		So(m.Video.GetVideoStream(), ShouldEqual, "/library/parts/2478/file.mp4")
		So(m.Video.GetSubtitleStream(), ShouldEqual, "/library/streams/5889")
	})
}
