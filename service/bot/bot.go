package bot

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime/debug"
	"time"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/utils"
	"github.com/ProtobufBot/Go-Mirai-Client/config"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/util"
	"github.com/fanliao/go-promise"
	log "github.com/sirupsen/logrus"
)

var Cli *client.QQClient
var LoginToken []byte

func InitDevice(uin int64) {
	// 默认 device/device-qq.json
	devicePath := path.Join("device", fmt.Sprintf("device-%d.json", uin))

	// 优先使用参数目录
	if config.Device != "" {
		devicePath = config.Device
	}

	deviceDir := path.Dir(devicePath)
	if !util.PathExists(deviceDir) {
		log.Infof("%+v 目录不存在，自动创建", deviceDir)
		if err := os.MkdirAll(deviceDir, 0777); err != nil {
			log.Warnf("failed to mkdir deviceDir, err: %+v", err)
		}
	}

	log.Info("生成随机设备信息")
	client.SystemDeviceInfo.AndroidId = []byte("MIRAI.123456.001")
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
	client.SystemDeviceInfo.WifiSSID = []byte("TP-LINK-" + utils.RandomStringRange(6, "ABCDEF1234567890"))
	client.SystemDeviceInfo.IpAddress = []byte{192, 168, 1, byte(100 + uin%100)}
	client.SystemDeviceInfo.Protocol = client.IPad
	client.SystemDeviceInfo.VendorOSName = []byte("gmc")

	if util.PathExists(devicePath) {
		log.Infof("使用 %s 内的设备信息覆盖设备信息", devicePath)
		if err := client.SystemDeviceInfo.ReadJson([]byte(util.ReadAllText(devicePath))); err != nil {
			util.FatalError(fmt.Errorf("failed to load device info, err: %+v", err))
		}
	}

	log.Infof("保存设备信息到文件 %s", devicePath)
	err := ioutil.WriteFile(devicePath, client.SystemDeviceInfo.ToJson(), 0644)
	if err != nil {
		log.Warnf("写设备信息文件 %s 失败", devicePath)
	}
}

func InitLog(cli *client.QQClient) {
	cli.OnLog(func(c *client.QQClient, e *client.LogEvent) {
		switch e.Type {
		case "INFO":
			log.Info("MiraiGo -> " + e.Message)
		case "ERROR":
			log.Error("MiraiGo -> " + e.Message)
			log.Errorf("%+v", string(debug.Stack()))
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
		if cli.Online {
			return
		}
		cli.Disconnect()
		var times = 1
		for {
			if cli.Online {
				log.Warn("Bot已登录")
				return
			}
			if times > retryCount {
				break
			}
			log.Warnf("Bot已离线 (%v)，将在 %v 秒后尝试重连. 重连次数：%v",
				e.Message, retryInterval, times)
			times++
			time.Sleep(time.Second * time.Duration(retryInterval))

			// 尝试token登录
			if err := cli.TokenLogin(LoginToken); err != nil {
				log.Errorf("failed to relogin with token, try to login with password, %+v", err)
				cli.Disconnect()
			} else {
				LoginToken = cli.GenToken()
				log.Info("succeed to relogin with token")
				return
			}

			time.Sleep(time.Second)

			// 尝试密码登录
			ok, err := Login(cli)

			if err != nil {
				log.Errorf("重连失败: %v", err)
				cli.Disconnect()
				continue
			}
			if ok {
				LoginToken = cli.GenToken()
				log.Info("重连成功")
				return
			}
		}
		util.FatalError(fmt.Errorf("failed to reconnect: 重连次数达到设置的上限值"))
	})
}
