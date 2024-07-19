package gmc_android

import (
	"os"

	"github.com/2mf8/Go-Lagrange-Client/pkg/config"
	"github.com/2mf8/Go-Lagrange-Client/pkg/gmc"

	log "github.com/sirupsen/logrus"
	_ "golang.org/x/mobile/bind"
)

// gomobile bind -target=android -androidapi=21 ./service/gmc_android
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
