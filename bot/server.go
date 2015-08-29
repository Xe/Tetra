package tetra

import (
	"strconv"

	"github.com/rcrowley/go-metrics"
)

// Struct Server holds information for a TS6 server.
type Server struct {
	Sid     string
	Name    string
	Gecos   string
	Links   []*Server
	Count   int
	Counter metrics.Gauge
	Hops    int
	Capab   []string
}

// AddClient increments the server client counter.
func (s *Server) AddClient() {
	s.Count++
}

// DelClient decrements the server client counter.
func (s *Server) DelClient() {
	s.Count--
}

// NewServer allocates a new server struct, fitting it into the network.
func NewServer(parent *Server, name, gecos, id, hops string) *Server {
	s := &Server{
		Sid:     id,
		Name:    name,
		Gecos:   gecos,
		Links:   []*Server{parent},
		Counter: metrics.NewGauge(),
	}

	s.Hops, _ = strconv.Atoi(hops)
	parent.Links = append(parent.Links, s)

	return s
}
