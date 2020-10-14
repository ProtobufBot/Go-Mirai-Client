package main

import (
	"bufio"
	"github.com/ProtobufBot/Go-Mirai-Client/service/bot"
	"os"
	"strconv"
	"time"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/plugin"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/util"
	"github.com/ProtobufBot/Go-Mirai-Client/service/plugins"
	log "github.com/sirupsen/logrus"
)

func main() {
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

	log.Infof("开始读取设备信息")
	bot.InitDevice()
	log.Infof("设备信息 %+v", client.SystemDeviceInfo)

	log.Infof("创建机器人 %+v", uin)
	cli := client.NewClient(uin, pass)

	log.Infof("初始化日志")
	bot.InitLog(cli)

	log.Infof("加载日志插件 Log")
	plugin.AddPrivateMessagePlugin(plugins.LogPrivateMessage)
	plugin.AddGroupMessagePlugin(plugins.LogGroupMessage)

	log.Infof("加载测试插件 Hello")
	plugin.AddPrivateMessagePlugin(plugins.HelloPrivateMessage)

	log.Infof("加载上报插件 Report")
	plugin.AddPrivateMessagePlugin(plugins.ReportPrivateMessage)
	plugin.AddGroupMessagePlugin(plugins.ReportGroupMessage)

	plugin.Serve(cli)
	log.Infof("插件加载完成")

	log.Infof("登录中...")
	ok, err := bot.Login(cli)
	if err != nil {
		log.Errorf("登录失败%v", err)
		time.Sleep(5 * time.Second)
		os.Exit(0)
		return
	}
	if ok {
		log.Infof("登录成功")
	} else {
		log.Infof("登录失败")
	}

	time.Sleep(5 * time.Second)

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

	log.Infof("刷新好友列表")
	util.Check(cli.ReloadFriendList())
	log.Infof("共加载 %v 个好友.", len(cli.FriendList))

	log.Infof("刷新群列表")
	util.Check(cli.ReloadGroupList())
	log.Infof("共加载 %v 个群.", len(cli.GroupList))

	bot.ConnectUniversal(cli)

	bot.SetRelogin(cli, 30, 30)
	_, _ = Console.ReadString('\n')

}
