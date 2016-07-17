package gdm

import (
	"champ/plex/model"
	"net"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
)

type GDM struct {
	regAddr         *net.UDPAddr
	dscAddr         *net.UDPAddr
	player          *model.Player
	port            string
	maxDatagramSize int
	conn            *net.UDPConn
	shutdown        bool
	shutdownCh      chan struct{}
	shutdownLock    sync.Mutex
}

func NewAdvertiser(player *model.Player, port string) (error, *GDM) {
	mip := net.ParseIP("239.0.0.250")
	s := &GDM{
		regAddr:         &net.UDPAddr{IP: mip, Port: 32413},
		dscAddr:         &net.UDPAddr{IP: mip, Port: 32412},
		player:          player,
		port:            port,
		maxDatagramSize: 8192,
		shutdownCh:      make(chan struct{}),
	}
	conn, err := net.ListenMulticastUDP("udp", nil, s.dscAddr)
	if err != nil {
		return err, nil
	}
	s.conn = conn
	conn.SetReadBuffer(s.maxDatagramSize)
	s.register()
	go s.recv()
	return nil, s
}

func (s *GDM) msgHandler(dst *net.UDPAddr, data string) {
	if strings.Contains(data, "M-SEARCH * HTTP/1.1") {
		res := newResponse()
		res.AddLine("HTTP/1.0 200 OK")
		res.ClientInfo(s.player, s.port)
		res.WriteTo(s.conn, dst)

	}
}

func (s *GDM) register() {
	res := newResponse()
	res.AddLine("HELLO * HTTP/1.0")
	res.ClientInfo(s.player, s.port)
	res.WriteTo(s.conn, s.regAddr)
}

func (s *GDM) unRegister() {
	res := newResponse()
	res.AddLine("BYE * HTTP/1.0")
	res.ClientInfo(s.player, s.port)
	res.WriteTo(s.conn, s.regAddr)
}

func (s *GDM) recv() error {
	for !s.shutdown {
		b := make([]byte, s.maxDatagramSize)
		n, src, err := s.conn.ReadFromUDP(b)
		if err != nil {
			log.Errorln(err)
			break
		} else {
			data := string(b[:n])
			s.msgHandler(src, data)
		}
	}

	return nil
}

// Shutdown is used to shutdown the listener
func (s *GDM) Shutdown() error {
	s.shutdownLock.Lock()
	defer s.shutdownLock.Unlock()

	if s.shutdown {
		return nil
	}
	s.shutdown = true
	close(s.shutdownCh)
	s.unRegister()
	s.conn.Close()
	return nil
}
