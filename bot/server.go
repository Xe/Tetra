package cod

type Server struct {
	Sid   string
	Name  string
	Gecos string
	Links []*Server
}
