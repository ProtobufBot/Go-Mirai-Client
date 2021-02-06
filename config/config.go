package config

import (
	"encoding/json"
	"github.com/pkg/errors"
)

var (
	Fragment = false  // 是否分片
	Conf     = &GmcConfig{
		SMS:  false,
		Port: "9000",
		ServerGroups: []*ServerGroup{
			{Name: "default", Disabled: false, Urls: []string{"ws://localhost:8081/ws/cq/"}},
		},
	}
)

type GmcConfig struct {
	Port         string         `json:"port"`
	SMS          bool           `json:"sms"`           // 设备锁是否优先使用短信认证
	ServerGroups []*ServerGroup `json:"server_groups"` // 服务器组
}

type ServerGroup struct {
	Name     string   `json:"name"`
	Disabled bool     `json:"disabled"`
	Urls     []string `json:"urls"`
	// TODO event filter, msg filter
}

func (g *GmcConfig) ReadJson(d []byte) error {
	var fileConfig GmcConfig
	if err := json.Unmarshal(d, &fileConfig); err != nil {
		return errors.Wrap(err, "failed to unmarshal json GmcConfig")
	}
	*g = fileConfig
	return nil
}

func (g *GmcConfig) ToJson() []byte {
	b, _ := json.MarshalIndent(g,"","    ")
	return b
}
