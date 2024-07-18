package config

import (
	"fmt"
	"os"
	"time"

	"github.com/2mf8/Go-Lagrange-Client/pkg/util"
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

type Setting struct {
	Platform   string
	AppVersion string
	SignServer string
}

var SettingPath = "setting"
var AllSetting *Setting = &Setting{}

func AllSettings() *Setting {
	_, err := toml.DecodeFile("setting/setting.toml", AllSetting)
	if err != nil {
		return AllSetting
	}
	return AllSetting
}

func ReadSetting() Setting {
	tomlData := `# linux / macos / windows, 默认linux
Platform = "linux"
# linux[3.1.2-13107,3.2.10-25765] macos[6.9.20-17153] windows[9.9.12-25493]
AppVersion = "3.1.2-13107"
# 默认 linux 3.1.2-13107 可用 master:https://sign.lagrangecore.org/api/sign,Mirror:https://sign.0w0.ing/api/sign
# linux 3.2.10-25765 暂不提供SignServer
SignServer = "https://sign.lagrangecore.org/api/sign"
	`
	if !util.PathExists(SettingPath) {
		if err := os.MkdirAll(SettingPath, 0777); err != nil {
			log.Warnf("failed to mkdir")
			return *AllSetting
		}
	}
	_, err := os.Stat(fmt.Sprintf("%s/setting.toml", SettingPath))
	if err != nil {
		_ = os.WriteFile(fmt.Sprintf("%s/setting.toml", SettingPath), []byte(tomlData), 0644)
		log.Warn("已生成配置文件 conf.toml ,请修改后重新启动程序。")
		log.Info("该程序将于5秒后退出！")
		time.Sleep(time.Second * 5)
		os.Exit(1)
	}
	AllSetting = AllSettings()
	fmt.Println(AllSetting.AppVersion, AllSetting.Platform, AllSetting.SignServer)
	return *AllSetting
}
