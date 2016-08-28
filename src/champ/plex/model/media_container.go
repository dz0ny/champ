package model

import "encoding/xml"

type MediaContainer struct {
	CommandID         string `xml:"commandID,attr,omitempty"`
	Identifier        string `xml:"identifier,attr,omitempty"`
	MachineIdentifier string `xml:"machineIdentifier,attr,omitempty"`
	Location          string `xml:"location,attr,omitempty"`
	Size              int    `xml:"size,attr,omitempty"`
	LibraryTitle      string `xml:"librarySectionTitle,attr,omitempty"`

	PlayQueueID                     int `xml:"playQueueID,attr,omitempty"`
	PlayQueueSelectedItemID         int `xml:"playQueueSelectedItemID,attr,omitempty"`
	PlayQueueSelectedItemOffset     int `xml:"playQueueSelectedItemOffset,attr,omitempty"`
	PlayQueueSelectedMetadataItemID int `xml:"playQueueSelectedMetadataItemID,attr,omitempty"`
	PlayQueueTotalCount             int `xml:"playQueueTotalCount,attr,omitempty"`

	Player    *Player    `xml:"Player"`
	Video     []Video    `xml:"Video"`
	Audio     []Audio    `xml:"Track"`
	Timelines []Timeline `xml:"Timeline"`
}

func (c *MediaContainer) String() string {
	output, err := xml.MarshalIndent(*c, "  ", "    ")
	if err != nil {
		return string(output)
	}
	return err.Error()
}
