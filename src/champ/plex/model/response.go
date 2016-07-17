package model

type Response struct {
	Code   string `xml:"code,attr"`
	Status string `xml:"status,attr"`
}
