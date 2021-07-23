package main

import (
	"fmt"

	"github.com/ProtobufBot/Go-Mirai-Client/pkg/config"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/gmc"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/util"

	log "github.com/sirupsen/logrus"
	"github.com/zserge/lorca"
)

func main() {
	gmc.Start()
	ui, err := lorca.New(fmt.Sprintf("http://localhost:%s", config.Port), "", 1024, 768)
	if err != nil {
		util.FatalError(err)
		return
	}
	defer ui.Close()
	<-ui.Done()
	log.Info("UI exit.")
}
