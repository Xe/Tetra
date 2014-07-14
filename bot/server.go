package tetra

import (
	"github.com/rcrowley/go-metrics"
)

type Server struct {
	Sid     string    `json:"sid"`
	Name    string    `json:"name"`
	Gecos   string    `json:"gecos"`
	Links   []*Server `json:"links"`
	count   int       `json:"usercount"`
	Counter metrics.Gauge
}

func (s *Server) AddClient() {
	s.count++
}

func (s *Server) DelClient() {
	s.count--
}
