package tetra

type Server struct {
	Sid   string
	Name  string
	Gecos string
	Links []*Server
}
