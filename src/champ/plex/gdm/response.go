package gdm

import (
	"bytes"
	"champ/plex/model"
	"fmt"
	"net"
)

type Response struct {
	buffer bytes.Buffer
}

func newResponse() *Response {
	return &Response{bytes.Buffer{}}
}

func (r *Response) AddHeader(key, value string) {
	r.AddLine(fmt.Sprintf("%s: %s", key, value))
}

func (r *Response) ClientInfo(player *model.Player, port string) {
	r.AddHeader("Name", player.Title)
	r.AddHeader("Port", port)
	r.AddHeader("Product", player.Product)
	r.AddHeader("Protocol", player.Protocol)
	r.AddHeader("Protocol-Version", player.ProtocolVersion)
	r.AddHeader("Protocol-Capabilities", player.ProtocolCapabilities)
	r.AddHeader("Resource-Identifier", player.MachineIdentifier)
	r.AddHeader("Device-Class", player.DeviceClass)
	r.AddHeader("Version", player.PlatformVersion)

	r.AddHeader("Content-Type", "plex/media-player")
}

func (r *Response) AddLine(line string) {
	r.buffer.WriteString(line + "\r\n")
}

func (r *Response) WriteTo(conn *net.UDPConn, addr *net.UDPAddr) {
	conn.WriteToUDP(r.buffer.Bytes(), addr)
	r.buffer.Reset()
}
