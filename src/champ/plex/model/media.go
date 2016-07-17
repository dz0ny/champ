package model

type Media struct {
	VideoResolution string `xml:"videoResolution,attr"`
	ID              string `xml:"id,attr"`
	VideoCodec      string `xml:"videoCodec,attr"`
	Duration        string `xml:"duration,attr"`
	Part            Part   `xml:"Part"`
}

type Part struct {
	Duration string   `xml:"duration,attr"`
	ID       string   `xml:"id,attr"`
	Key      string   `xml:"key,attr"`
	File     string   `xml:"file,attr"`
	Stream   []Stream `xml:"Stream"`
}

type Stream struct {
	ID       string `xml:"id,attr"`
	Key      string `xml:"key,attr"`
	Codec    string `xml:"codec,attr"`
	Format   string `xml:"format,attr"`
	Duration string `xml:"duration,attr"`
}
