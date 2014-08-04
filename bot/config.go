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

type StatsConfig struct {
	Host     string `json:"host"`
	Database string `json:"db"`
	Username string `json:"user"`
	Password string `json:"pass"`
}

type Config struct {
	Autoload []string          `json:"autoload"`
	Services []*ServiceConfig  `json:"services"`
	Server   *ServerConfig     `json:"myinfo"`
	Uplink   *UplinkConfig     `json:"uplink"`
	Stats    *StatsConfig      `json:"stats"`
	ApiKeys  map[string]string `json:"apikeys"`
}

func NewConfig(fname string) (conf *Config, err error) {
	contents, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(contents, &conf)

	return
}
