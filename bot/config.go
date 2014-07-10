package tetra

import (
	"encoding/json"
	"io/ioutil"
)

type ServerConfig struct {
	Name  string `json:"name"`
	Gecos string `json:"gecos"`
}

type ServiceConfig struct {
	Nick  string `json:"nick"`
	User  string `json:"user"`
	Host  string `json:"host"`
	Gecos string `json:"gecos"`
	Name  string `json:"name"`
}

type UplinkConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	Ssl      bool   `json:"ssl"`
}

type Config struct {
	Autoload []string         `json:"autoload"`
	Services []*ServiceConfig `json:"services"`
	Server   *ServerConfig    `json:"myinfo"`
	RRDPath  string           `json:"rrdpath"`
}

func NewConfig(fname string) (conf *Config) {
	contents, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(contents, conf)

	return
}
