package model

type Player struct {
	Title                string `xml:"title,attr"`
	Protocol             string `xml:"protocol,attr"`
	ProtocolVersion      string `xml:"protocolVersion,attr"`
	ProtocolCapabilities string `xml:"protocolCapabilities,attr"`
	MachineIdentifier    string `xml:"machineIdentifier,attr"`
	Product              string `xml:"product,attr"`
	Platform             string `xml:"platform,attr"`
	PlatformVersion      string `xml:"platformVersion,attr"`
	DeviceClass          string `xml:"deviceClass,attr"`
	LastCommandID        string `xml:"lastCommandID,attr"`
}
