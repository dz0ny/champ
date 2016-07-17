package plex

import (
	"bytes"
	"champ/plex/model"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

type PlexClient struct {
	Remote *Remote
}

type Remote struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func NewRemote(scheme, address, port string) *Remote {
	// We need extra low latencies
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 1 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 1 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 1,
		Transport: netTransport,
	}
	return &Remote{&url.URL{Scheme: scheme, Host: fmt.Sprintf("%s:%s", address, port)}, netClient}
}

func (r *Remote) Path(path string) string {
	pathURL, pathErr := url.Parse(path)
	if pathErr == nil {
		return r.baseURL.ResolveReference(pathURL).String()
	}
	return ""
}
func (r *Remote) String() string {
	return r.baseURL.String()
}

func (r *Remote) Notify(data model.MediaContainer) {
	output, err := xml.MarshalIndent(data, "  ", "    ")
	if err != nil {
		log.Println(err)
		return
	}

	req, err := http.NewRequest("POST", r.Path("/:/timeline"), bytes.NewReader(output))
	req.Header.Set("X-Plex-Client-Identifier", data.MachineIdentifier)
	r.httpClient.Do(req)
}

func NewPlexClient(remote *Remote, httpClient *http.Client) *PlexClient {
	pc := &PlexClient{
		Remote: remote,
	}
	return pc
}

func (pc *PlexClient) GetMedia(url string) (error, model.MediaContainer) {

	var media model.MediaContainer
	resp, err := pc.Remote.httpClient.Get(pc.Remote.Path(url))
	defer resp.Body.Close()
	if err != nil {
		return err, media
	}

	if code := resp.StatusCode; 200 <= code && code <= 299 {
		err = xml.NewDecoder(resp.Body).Decode(&media)
	} else {
		return errors.New("Failed getting resource"), media
	}
	if err != nil {
		return err, media
	}
	return nil, media
}
