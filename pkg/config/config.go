package config

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
