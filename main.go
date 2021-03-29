package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/ProtobufBot/Go-Mirai-Client/config"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/util"
	"github.com/ProtobufBot/Go-Mirai-Client/service/handler"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	customFormatter := &log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceColors:     true,
	}
	log.SetFormatter(customFormatter)
}

func main() {
	gmcConfigPath := "gmc_config.json"
	LoadGmcConfigFile(gmcConfigPath)  // 如果文件存在，从文件读取gmc config
	LoadEnvConfig()                   // 如果环境变量存在，从环境变量读取gmc config，并覆盖
	WriteGmcConfigFile(gmcConfigPath) // 内存中的gmc config写到文件
	log.Infof("gmc config: %+v", util.MustMarshal(config.Conf))

	CreateBotIfEnvAccountExist() // 如果环境变量存在，使用环境变量创建机器人 UIN PASSWORD
	InitGin()                    // 初始化GIN HTTP管理

	select {}
}

func LoadGmcConfigFile(filePath string) {
	if util.PathExists(filePath) {
		if err := config.Conf.ReadJson([]byte(util.ReadAllText(filePath))); err != nil {
			log.Errorf("failed to read gmc config file %s, %+v", filePath, err)
		}
	}
}

func LoadEnvConfig() {
	if os.Getenv("SMS") == "1" {
		config.Conf.SMS = true
	}
	envWsUrl := os.Getenv("WS_URL")
	if envWsUrl != "" {
		config.Conf.ServerGroups = []*config.ServerGroup{
			{Name: "default", Urls: []string{envWsUrl}},
		}
	}
	envPort := os.Getenv("PORT")
	if envPort != "" {
		config.Conf.Port = envPort
	}
}

func WriteGmcConfigFile(filePath string) {
	if err := ioutil.WriteFile(filePath, config.Conf.ToJson(), 0644); err != nil {
		log.Warnf("failed to write gmc config file %s, %+v", filePath, err)
	}
}

func CreateBotIfEnvAccountExist() {
	envUin := os.Getenv("UIN")
	envPass := os.Getenv("PASSWORD")
	if envUin != "" || envPass != "" {
		uin, err := strconv.ParseInt(envUin, 10, 64)
		if err != nil {
			log.Errorf("环境变量账号错误")
		}
		log.Infof("使用环境变量创建机器人 %d", uin)
		go func() {
			handler.CreateBotImpl(uin, envPass)
		}()
	}
}

func InitGin() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(handler.CORSMiddleware())
	router.Static("/", "./static")
	router.POST("/bot/create/v1", handler.CreateBot)
	router.POST("/bot/list/v1", handler.ListBot)
	router.POST("/captcha/list/v1", handler.ListCaptcha)
	router.POST("/captcha/solve/v1", handler.SolveCaptcha)
	router.POST("/qrcode/fetch/v1", handler.FetchQrCode)
	router.POST("/qrcode/query/v1", handler.QueryQRCodeStatus)
	realPort, err := RunGin(router, ":"+config.Conf.Port)
	if err != nil {
		util.FatalError(fmt.Errorf("failed to run gin, err: %+v", err))
	}
	config.Conf.Port = realPort
	log.Infof("端口号 %s", realPort)
	log.Infof(fmt.Sprintf("浏览器打开 http://localhost:%s/ 设置机器人", realPort))
}

func RunGin(engine *gin.Engine, port string) (string, error) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return "", err
	}
	_, randPort, _ := net.SplitHostPort(ln.Addr().String())
	go func() {
		if err := http.Serve(ln, engine); err != nil {
			util.FatalError(fmt.Errorf("failed to serve http, err: %+v", err))
		}
	}()
	return randPort, nil
}
