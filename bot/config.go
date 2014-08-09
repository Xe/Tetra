package tetra

import (
	"io/ioutil"

	"gopkg.in/yaml.v1"
)

type ServerConfig struct {
	Name      string
	Gecos     string
	Sid       string
	StaffChan string
	SnoopChan string
}

type ServiceConfig struct {
	Nick  string
	User  string
	Host  string
	Gecos string
	Name  string
}

type UplinkConfig struct {
	Host     string
	Port     string
	Password string
	Ssl      bool
}

type StatsConfig struct {
	Host     string
	Database string
	Username string
	Password string
}

type Config struct {
	Autoload []string
	Services []*ServiceConfig
	Server   *ServerConfig
	Uplink   *UplinkConfig
	Stats    *StatsConfig
	ApiKeys  map[string]string
}

func NewConfig(fname string) (conf *Config, err error) {
	contents, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(contents, &conf)

	return
}
