package gmc_android

import (
	"os"

	"github.com/ProtobufBot/Go-Mirai-Client/pkg/config"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/gmc"

	log "github.com/sirupsen/logrus"
	_ "golang.org/x/mobile/bind"
)

var logger AndroidLogger

// SetSms 设置是否短信优先
func SetSms(sms bool) {
	config.SMS = sms
}

func Chdir(dir string) {
	_ = os.Chdir(dir)
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
