package tetra

import (
	"github.com/rcrowley/go-metrics"
)

// Struct Server holds information for a TS6 server.
type Server struct {
	Sid     string
	Name    string
	Gecos   string
	Links   []*Server
	count   int
	Counter metrics.Gauge
	Hops    int
	Capab   []string
}

// AddClient increments the server client counter.
func (s *Server) AddClient() {
	s.count++
}

// DelClient decrements the server client counter.
func (s *Server) DelClient() {
	s.count--
}
