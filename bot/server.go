package tetra

import (
	"github.com/rcrowley/go-metrics"
)

type Server struct {
	Sid     string
	Name    string
	Gecos   string
	Links   []*Server
	count   int
	Counter metrics.Gauge
}

func (s *Server) AddClient() {
	s.count++
}

func (s *Server) DelClient() {
	s.count--
}
