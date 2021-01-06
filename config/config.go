package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var FRAGMENT = false

type GmcConfig struct {
	Server Server `yaml:"server"`
	Bot    Bot    `yaml:"bot"`
}

type Server struct {
	Port int32 `yaml:"port"`
}

type Bot struct {
	Client Client `yaml:"client"`
}

type Client struct {
	WsUrl string `yaml:"ws-url"`
}

func LoadConfig(path string) (*GmcConfig, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config GmcConfig
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
