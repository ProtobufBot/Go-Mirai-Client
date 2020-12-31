package bot

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/utils"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/util"
	"github.com/fanliao/go-promise"
	log "github.com/sirupsen/logrus"
)

var Cli *client.QQClient

func InitDevice(path string) {
	if !util.PathExists(path) {
		log.Warn("虚拟设备信息不存在, 将自动生成随机设备.")
		client.GenRandomDevice()
		client.SystemDeviceInfo.Display = []byte("GMC." + utils.RandomStringRange(6, "0123456789") + ".001")
		client.SystemDeviceInfo.FingerPrint = []byte("pbbot/gmc/gmc:10/PBBOT.200324.001/" + utils.RandomStringRange(7, "0123456789") + ":user/release-keys")
		client.SystemDeviceInfo.ProcVersion = []byte("Linux 5.4.0-54-generic" + utils.RandomString(8) + " (android-build@gmail.com)")
		client.SystemDeviceInfo.AndroidId = client.SystemDeviceInfo.Display
		client.SystemDeviceInfo.Device = []byte("gmc")
		client.SystemDeviceInfo.Board = []byte("gmc")
		client.SystemDeviceInfo.Model = []byte("gmc")
		client.SystemDeviceInfo.Brand = []byte("pbbot")
		client.SystemDeviceInfo.Product = []byte("gmc")
		client.SystemDeviceInfo.Protocol = client.MacOS
		_ = ioutil.WriteFile(path, client.SystemDeviceInfo.ToJson(), 0644)
		log.Infof("已生成设备信息并保存到 %s 文件.", path)
	}
	log.Infof("将使用 %s 内的设备信息运行", path)
	client.SystemDeviceInfo.IpAddress = []byte{192, 168, 31, 101}
	if err := client.SystemDeviceInfo.ReadJson([]byte(util.ReadAllText(path))); err != nil {
		log.Errorf("加载设备信息失败: %v", err)
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}
}

func InitLog(cli *client.QQClient) {
	cli.OnLog(func(c *client.QQClient, e *client.LogEvent) {
		switch e.Type {
		case "INFO":
			log.Info("MiraiGo -> " + e.Message)
		case "ERROR":
			log.Error("MiraiGo -> " + e.Message)
		case "DEBUG":
			log.Debug("MiraiGo -> " + e.Message)
		}
	})

	cli.OnServerUpdated(func(bot *client.QQClient, e *client.ServerUpdatedEvent) bool {
		log.Infof("收到服务器地址更新通知, 将在下一次重连时应用. ")
		return true // 如果是 false 表示不应用
	})
}

func Login(cli *client.QQClient) (bool, error) {
	cli.AllowSlider = true
	rsp, err := cli.Login()
	if err != nil {
		return false, err
	}

	v, err := promise.Start(func() bool {
		ok, err := ProcessLoginRsp(cli, rsp)
		if err != nil {
			log.Errorf("登陆遇到错误2:%v", err)
			time.Sleep(5 * time.Second)
			os.Exit(0)
		}
		return ok
	}()).Get()
	if err != nil {
		return false, err
	}
	return v.(bool), nil
}

func SetRelogin(cli *client.QQClient, retryInterval int, retryCount int) {
	cli.OnDisconnected(func(bot *client.QQClient, e *client.ClientDisconnectedEvent) {
		var times = 1
		for {
			if cli.Online {
				log.Warn("Bot已登录")
				return
			}
			if retryCount == 0 {
			} else if times > retryCount {
				break
			}
			log.Warnf("Bot已离线 (%v)，将在 %v 秒后尝试重连. 重连次数：%v",
				e.Message, retryInterval, times)
			times++
			time.Sleep(time.Second * time.Duration(retryInterval))
			ok, err := Login(cli)

			if err != nil {
				log.Errorf("重连失败: %v", err)
				continue
			}
			if ok {
				log.Info("重连成功")
				return
			}
		}
		log.Errorf("重连失败: 重连次数达到设置的上限值")
		time.Sleep(5 * time.Second)
		os.Exit(0)
	})
}
