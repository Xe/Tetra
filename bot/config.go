package tetra

import (
	"encoding/json"
	"io/ioutil"
)

type ServerConfig struct {
	Name      string `json:"name"`
	Gecos     string `json:"gecos"`
	Sid       string `json:"sid"`
	StaffChan string `json:"staffchan"`
	SnoopChan string `json:"snoopchan"`
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
	Uplink   *UplinkConfig    `json:"uplink"`
	RRDPath  string           `json:"rrdpath"`
}

func NewConfig(fname string) (conf *Config, err error) {
	contents, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(contents, &conf)

	return
}
