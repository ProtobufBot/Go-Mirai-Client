package bot

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/ProtobufBot/Go-Mirai-Client/pkg/util"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/fanliao/go-promise"
	log "github.com/sirupsen/logrus"
)

var (
	// TODO sync
	Clients     = map[int64]*client.QQClient{}
	LoginTokens = map[int64][]byte{}
)

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
	LoginTokens[cli.Uin] = cli.GenToken()
	cli.OnDisconnected(func(bot *client.QQClient, e *client.ClientDisconnectedEvent) {
		if bot.Online {
			return
		}
		bot.Disconnect()
		var times = 1
		for IsClientExist(bot.Uin) {
			if bot.Online {
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
			if err := bot.TokenLogin(LoginTokens[bot.Uin]); err != nil {
				log.Errorf("failed to relogin with token, try to login with password, %+v", err)
				bot.Disconnect()
			} else {
				LoginTokens[cli.Uin] = bot.GenToken()
				log.Info("succeed to relogin with token")
				return
			}

			time.Sleep(time.Second)

			// 尝试密码登录
			ok, err := Login(bot)

			if err != nil {
				log.Errorf("重连失败: %v", err)
				bot.Disconnect()
				continue
			}
			if ok {
				LoginTokens[bot.Uin] = bot.GenToken()
				log.Info("重连成功")
				return
			}
		}
		log.Errorf("failed to reconnect: 重连次数达到设置的上限值, %+v", cli.Uin)
		ReleaseClient(cli)
	})
}

// ReleaseClient 断开连接并释放资源
func ReleaseClient(cli *client.QQClient) {
	cli.Release()
	delete(Clients, cli.Uin)
	delete(LoginTokens, cli.Uin)
}

func IsClientExist(uin int64) bool {
	_, ok := Clients[uin]
	return ok
}
