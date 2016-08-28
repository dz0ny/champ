package model

type Video struct {
	Key             string `xml:"key,attr"`
	PlayQueueItemID int    `xml:"playQueueItemID,attr"`
	RatingKey       string `xml:"ratingKey,attr"`
	Type            string `xml:"type,attr"`
	Title           string `xml:"title,attr"`
	Thumbnail       string `xml:"thumb,attr"`
	Summary         string `xml:"summary,attr"`
	Art             string `xml:"art,attr"`
	Duration        string `xml:"duration,attr"`
	Guid            string `xml:"guid,attr"`
	Media           Media  `xml:"Media"`
}

func (c *Video) VideoStream() string {
	return c.Media.Part.Key
}

func (c *Video) SubtitleStream() (string, string) {
	for _, el := range c.Media.Part.Stream {
		if el.Format == "srt" {
			return el.ID, el.Key
		}
	}
	return "", ""
}
