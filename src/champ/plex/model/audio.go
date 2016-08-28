package model

type Audio struct {
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

func (a *Audio) AudioStream() (string, string) {
	return a.Media.Part.ID, a.Media.Part.Key
}
