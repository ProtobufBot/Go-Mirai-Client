package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/ProtobufBot/Go-Mirai-Client/pkg/util"

	log "github.com/sirupsen/logrus"
)

//go:generate go run github.com/a8m/syncmap -o "gen_plugin_map.go" -pkg config -name PluginMap "map[string]*Plugin"
var (
	Fragment = false // 是否分片
	Port     = "9000"
	SMS      = false
	Device   = ""
	Plugins  = &PluginMap{}
	HttpAuth = map[string]string{}
)

func init() {
	Plugins.Store("default", &Plugin{
		Name:         "default",
		Disabled:     false,
		Json:         false,
		Urls:         []string{"ws://localhost:8081/ws/cq/"},
		EventFilter:  []int32{},
		ApiFilter:    []int32{},
		RegexFilter:  "",
		RegexReplace: "",
		ExtraHeader: map[string][]string{
			"User-Agent": {"GMC"},
		},
	})
}

func ClearPlugins(pluginMap *PluginMap) {
	pluginMap.Range(func(key string, value *Plugin) bool {
		pluginMap.Delete(key)
		return true
	})
}

type Plugin struct {
	Name         string              `json:"-"`             // 功能名称
	Disabled     bool                `json:"disabled"`      // 不填false默认启用
	Json         bool                `json:"json"`          // json上报
	Urls         []string            `json:"urls"`          // 服务器列表
	EventFilter  []int32             `json:"event_filter"`  // 事件过滤
	ApiFilter    []int32             `json:"api_filter"`    // API过滤
	RegexFilter  string              `json:"regex_filter"`  // 正则过滤
	RegexReplace string              `json:"regex_replace"` // 正则替换
	ExtraHeader  map[string][]string `json:"extra_header"`  // 自定义请求头
	// TODO event filter, msg filter, regex filter, prefix filter, suffix filter
}

var PluginPath = "plugins"

func LoadPlugins() {
	if !util.PathExists(PluginPath) {
		return
	}
	files, err := ioutil.ReadDir(PluginPath)
	if err != nil {
		log.Warnf("failed to read plugin dir: %s", err)
		return
	}

	if len(files) == 0 {
		log.Warnf("plugin dir is empty")
		return
	}

	ClearPlugins(Plugins)
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		pluginName := strings.TrimSuffix(file.Name(), ".json")
		filepath := path.Join(PluginPath, file.Name())
		b, err := os.ReadFile(filepath)
		if err != nil {
			log.Warnf("failed to read plugin file: %s %s", filepath, err)
			continue
		}
		plugin := &Plugin{}
		if err := json.NewDecoder(bytes.NewReader(b)).Decode(plugin); err != nil {
			log.Warnf("failed to decode plugin file: %s %s", filepath, err)
			continue
		}
		plugin.Name = pluginName
		Plugins.Store(plugin.Name, plugin)
	}
}

func WritePlugins() {
	if !util.PathExists(PluginPath) {
		if err := os.MkdirAll(PluginPath, 0777); err != nil {
			log.Warnf("failed to mkdir")
			return
		}
	}
	DeletePluginFiles()
	Plugins.Range(func(key string, plugin *Plugin) bool {
		pluginFilename := fmt.Sprintf("%s.json", plugin.Name)
		filepath := path.Join(PluginPath, pluginFilename)
		b, err := json.MarshalIndent(plugin, "", "    ")
		if err != nil {
			log.Warnf("failed to marshal plugin, %s", plugin.Name)
			return true
		}
		if err := os.WriteFile(filepath, b, 0777); err != nil {
			log.Warnf("failed to write file, %s", pluginFilename)
			return true
		}
		return true
	})
}

func DeletePluginFiles() {
	files, err := ioutil.ReadDir(PluginPath)
	if err != nil {
		log.Warnf("failed to read plugin dir: %s", err)
	}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		filepath := path.Join(PluginPath, file.Name())
		if err := os.Remove(filepath); err != nil {
			log.Warnf("failed to remove plugin file: %s", filepath)
			continue
		}
	}
}
