package gmc_android

import (
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/config"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/device"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/gmc"

	log "github.com/sirupsen/logrus"
	_ "golang.org/x/mobile/bind"
)

var logger AndroidLogger

// SetPluginPath 设置插件配置路径
func SetPluginPath(pluginPath string) {
	config.PluginPath = pluginPath
}

// SetSms 设置是否短信优先
func SetSms(sms bool) {
	config.SMS = sms
}

// SetLogPath 设置日志目录
func SetLogPath(logPath string) {
	gmc.LogPath = logPath
}

// SetDevicePath 设置设备信息路径
func SetDevicePath(devicePath string) {
	device.Path = devicePath
}

// Start 启动主程序
func Start() {
	gmc.Start()
}

// SetLogger 设置日志输出
func SetLogger(androidLogger AndroidLogger) {
	logger = androidLogger
	log.SetOutput(&AndroidWriter{})
	log.SetFormatter(&AndroidFormatter{})
}
