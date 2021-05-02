package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ProtobufBot/Go-Mirai-Client/config"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/util"
	"github.com/ProtobufBot/Go-Mirai-Client/service/handler"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	sms          = false // 参数优先使用短信验证
	wsUrls       = ""    // websocket url
	port         = 9000  // 端口号
	uin    int64 = 0     // qq
	pass         = ""    //password
	device       = ""    // device file path
	help         = false // help
)

func init() {
	flag.BoolVar(&sms, "sms", false, "use sms captcha")
	flag.StringVar(&wsUrls, "ws_url", "", "websocket url")
	flag.IntVar(&port, "port", 9000, "admin http api port, 0 is random")
	flag.Int64Var(&uin, "uin", 0, "bot's qq")
	flag.StringVar(&pass, "pass", "", "bot's password")
	flag.StringVar(&device, "device", "", "device file")
	flag.BoolVar(&help, "help", false, "this help")
	flag.Parse()

	customFormatter := &log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceColors:     true,
	}
	log.SetFormatter(customFormatter)
	log.SetOutput(os.Stdout)
}

func main() {
	if help {
		flag.Usage()
		os.Exit(0)
	}

	gmcConfigPath := "gmc_config.json"
	LoadGmcConfigFile(gmcConfigPath)  // 如果文件存在，从文件读取gmc config
	LoadParamConfig()                 // 如果环境变量存在，从环境变量读取gmc config，并覆盖
	WriteGmcConfigFile(gmcConfigPath) // 内存中的gmc config写到文件
	log.Infof("gmc config: %+v", util.MustMarshal(config.Conf))

	CreateBotIfParamExist() // 如果环境变量存在，使用环境变量创建机器人 UIN PASSWORD
	InitGin()               // 初始化GIN HTTP管理

	select {}
}

func LoadGmcConfigFile(filePath string) {
	if util.PathExists(filePath) {
		if err := config.Conf.ReadJson([]byte(util.ReadAllText(filePath))); err != nil {
			log.Errorf("failed to read gmc config file %s, %+v", filePath, err)
		}
	}
}

func LoadParamConfig() {
	// sms是true，如果本来是true，不变。如果本来是false，变true
	if sms {
		config.SMS = true
	}

	if wsUrls != "" {
		wsUrlList := strings.Split(wsUrls, ",")
		config.Conf.ServerGroups = []*config.ServerGroup{}
		for i, wsUrl := range wsUrlList {
			config.Conf.ServerGroups = append(config.Conf.ServerGroups, &config.ServerGroup{Name: strconv.Itoa(i), Urls: []string{wsUrl}})
		}
	}

	if port != 9000 {
		config.Port = strconv.Itoa(port)
	}

	if device != "" {
		config.Device = device
	}
}

func WriteGmcConfigFile(filePath string) {
	if err := ioutil.WriteFile(filePath, config.Conf.ToJson(), 0644); err != nil {
		log.Warnf("failed to write gmc config file %s, %+v", filePath, err)
	}
}

func CreateBotIfParamExist() {
	if uin != 0 && pass != "" {
		log.Infof("使用参数创建机器人 %d", uin)
		go func() {
			handler.CreateBotImpl(uin, pass)
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
	realPort, err := RunGin(router, ":"+config.Port)
	if err != nil {
		util.FatalError(fmt.Errorf("failed to run gin, err: %+v", err))
	}
	config.Port = realPort
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
