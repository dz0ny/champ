package model

import "encoding/xml"

type MediaContainer struct {
	CommandID         string     `xml:"commandID,attr,omitempty"`
	Identifier        string     `xml:"identifier,attr,omitempty"`
	MachineIdentifier string     `xml:"machineIdentifier,attr,omitempty"`
	Location          string     `xml:"location,attr,omitempty"`
	LibraryTitle      string     `xml:"librarySectionTitle,attr,omitempty"`
	Player            *Player    `xml:"Player"`
	Video             *Video     `xml:"Video"`
	Audio             *Audio     `xml:"Track"`
	Timelines         []Timeline `xml:"Timeline"`
}

func (v *MediaContainer) VideoStream() string {
	if v.Video == nil {
		return ""
	}
	return v.Video.Media.Part.Key
}

func (v *MediaContainer) AudioStream() string {
	if v.Audio == nil {
		return ""
	}
	return v.Audio.Media.Part.Key
}

func (v *MediaContainer) SubtitleStream() string {
	for _, el := range v.Video.Media.Part.Stream {
		if el.Format == "srt" {
			return el.Key
		}
	}
	return ""
}

func (v *MediaContainer) String() string {
	output, err := xml.MarshalIndent(*v, "  ", "    ")
	if err != nil {
		return string(output)
	}
	return err.Error()
}
