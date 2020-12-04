package main

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"time"

	"github.com/ProtobufBot/Go-Mirai-Client/service/bot"
	"github.com/ProtobufBot/Go-Mirai-Client/service/handler"
	log "github.com/sirupsen/logrus"
)

func init() {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
}

func main() {
	port := "9000"
	conf, err := LoadConfig("application.yml")
	if err == nil && conf != nil {
		if conf.Bot.Client.WsUrl != "" {
			bot.WsUrl = conf.Bot.Client.WsUrl
		}
		if conf.Server.Port != 0 {
			port = strconv.Itoa(int(conf.Server.Port))
		}
	}
	envPort := os.Getenv("PORT")
	if envPort != "" {
		port = envPort
	}
	envWsUrl := os.Getenv("WS_URL")
	if envWsUrl != "" {
		bot.WsUrl = envWsUrl
	}
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

	log.Infof("端口号 %s", port)
	port = ":" + port
	router := gin.Default()
	router.Use(handler.CORSMiddleware())
	router.Static("/", "./static")
	router.POST("/bot/create/v1", handler.CreateBot)
	router.POST("/bot/list/v1", handler.ListBot)
	router.POST("/captcha/list/v1", handler.ListCaptcha)
	router.POST("/captcha/solve/v1", handler.SolveCaptcha)
	err = router.Run(port)
	if err != nil {
		log.Errorf("run server error %v", err)
	}
}

func TestBot() {
	Console := bufio.NewReader(os.Stdin)
	uinStr := os.Getenv("uin")
	pass := os.Getenv("pass")
	if uinStr == "" || pass == "" {
		log.Warnf("请在环境变量设置 uin 和 pass")
		time.Sleep(5 * time.Second)
		return
	}

	uin, err := strconv.ParseInt(uinStr, 10, 64)
	if err != nil {
		log.Warnf("uin 错误")
		time.Sleep(5 * time.Second)
		panic(err)
	}

	go func() {
		handler.CreateBotImpl(uin, pass)
	}()

	// TODO 改成 gin 处理验证码
	for {
		if bot.Captcha == nil {
			break
		}
		log.Infof("请输入验证码%+v", bot.Captcha)
		text, _ := Console.ReadString('\n')
		log.Infof("你输入的是:%v", text)
		err := bot.CaptchaPromise.Resolve(text)
		if err != nil {
			break
		}
		time.Sleep(5 * time.Second)
	}

	_, _ = Console.ReadString('\n')

}
