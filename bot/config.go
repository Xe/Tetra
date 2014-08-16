package tetra

import (
	"io/ioutil"

	"gopkg.in/yaml.v1"
)

// Atheme configuration
type AthemeConfig struct {
	URL      string
	Username string
	Password string
}

// Struct ServerConfig defines the server information for Tetra.
type ServerConfig struct {
	Name  string
	Gecos string
	Sid   string
}

// Struct GeneralConfig defines general configuration for Tetra.
type GeneralConfig struct {
	StaffChan string
	SnoopChan string
	Prefix    string
	Debug     bool
	Workers   int
}

// Struct ServiceConfig defines the configuration for a service.
type ServiceConfig struct {
	Nick   string
	User   string
	Host   string
	Gecos  string
	Name   string
	Certfp string
}

// Struct UplinkConfig defines the configuration of Tetra's uplink.
type UplinkConfig struct {
	Host     string
	Port     string
	Password string
	Ssl      bool
}

// Struct StatsConfig defines the InfluxxDB information for Tetra.
type StatsConfig struct {
	Host     string
	Database string
	Username string
	Password string
}

// Struct Config defines the configuration for Tetra.
type Config struct {
	Autoload []string
	Services []*ServiceConfig
	Server   *ServerConfig
	Uplink   *UplinkConfig
	Stats    *StatsConfig
	ApiKeys  map[string]string
	General  *GeneralConfig
	Modules  map[string]interface{}
	Etcd     struct {
		Machines []string
	}
	Atheme *AthemeConfig
}

// NewConfig returns a new Config instance seeded by the file at fname.
func NewConfig(fname string) (conf *Config, err error) {
	contents, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(contents, &conf)

	return
}
