package gmc_android

import (
	"os"

	"github.com/ProtobufBot/Go-Mirai-Client/pkg/config"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/gmc"

	log "github.com/sirupsen/logrus"
	_ "golang.org/x/mobile/bind"
)

var logger AndroidLogger

func SetSms(sms bool) {
	config.SMS = sms
}

func Chdir(dir string) {
	_ = os.Chdir(dir)
}

func Start() {
	gmc.Start()
}

func SetLogger(androidLogger AndroidLogger) {
	logger = androidLogger
	log.SetOutput(&AndroidWriter{})
	log.SetFormatter(&AndroidFormatter{})
}
