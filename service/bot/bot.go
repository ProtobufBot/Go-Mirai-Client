package bot

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/utils"
	"io/ioutil"
	"time"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/util"
	"github.com/fanliao/go-promise"
	log "github.com/sirupsen/logrus"
)

var Cli *client.QQClient

func InitDevice(path string) {
	log.Info("生成随机设备信息")
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
	client.SystemDeviceInfo.Protocol = client.IPad
	client.SystemDeviceInfo.IpAddress = []byte{192, 168, 31, 101}

	if util.PathExists(path) {
		log.Infof("使用 %s 内的设备信息覆盖随机设备信息", path)
		if err := client.SystemDeviceInfo.ReadJson([]byte(util.ReadAllText(path))); err != nil {
			util.FatalError(fmt.Errorf("failed to load device info, err: %+v", err))
		}
	}

	log.Infof("保存设备信息到文件 %s", path)
	err := ioutil.WriteFile(path, client.SystemDeviceInfo.ToJson(), 0644)
	if err != nil {
		log.Warnf("写设备信息文件 %s 失败", path)
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
			util.FatalError(fmt.Errorf("failed to login: %+v", err))
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
		util.FatalError(fmt.Errorf("failed to reconnect: 重连次数达到设置的上限值"))
	})
}
