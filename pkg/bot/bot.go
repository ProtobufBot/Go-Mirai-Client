package bot

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/wrapper"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/download"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

//go:generate go run github.com/a8m/syncmap -o "gen_client_map.go" -pkg bot -name ClientMap "map[int64]*client.QQClient"
//go:generate go run github.com/a8m/syncmap -o "gen_token_map.go" -pkg bot -name TokenMap "map[int64][]byte"
var (
	Clients     ClientMap
	LoginTokens TokenMap
)

type Logger struct {
}

type GMCLogin struct {
	DeviceSeed     int64
	ClientProtocol int32
	SignServer     string
}

var GTL *GMCLogin

func GmcTokenLogin() (g GMCLogin, err error) {
	_, err = toml.DecodeFile("deviceInfo.toml", &GTL)
	return *GTL, err
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || errors.Is(err, os.ErrExist)
}

func (l *Logger) Info(format string, args ...any) {
	log.Infof(format, args)
}
func (l *Logger) Warning(format string, args ...any) {
	log.Warnf(format, args)
}
func (l *Logger) Error(format string, args ...any) {
	log.Errorf(format, args)
}
func (l *Logger) Debug(format string, args ...any) {
	log.Debugf(format, args)
}
func (l *Logger) Dump(dumped []byte, format string, args ...any) {
}

func InitLog(cli *client.QQClient) {
	cli.SetLogger(&Logger{})

	cli.OnServerUpdated(func(bot *client.QQClient, e *client.ServerUpdatedEvent) bool {
		log.Infof("收到服务器地址更新通知, 将在下一次重连时应用. ")
		return true // 如果是 false 表示不应用
	})
}

func Login(cli *client.QQClient) (bool, error) {
	cli.AllowSlider = true
	if GTL.ClientProtocol == 1 && GTL.SignServer != "" {
		wrapper.DandelionEnergy = Energy
		wrapper.FekitGetSign = Sign
	} else if GTL.SignServer != "" {
		fmt.Println("SignServer 不支持该协议")
	}
	rsp, err := cli.Login()
	if rsp.Code == byte(45) && GTL.SignServer == "" {
		fmt.Println("您的账号被限制登录，请配置 SignServer 后重试")
	}
	if rsp.Code == byte(235) {
		fmt.Println("设备信息被封禁，请删除设备（device）文件夹里对应设备文件后重试")
	}
	if rsp.Code == byte(237) {
		fmt.Println("登录过于频繁，请在手机QQ登录并根据提示完成认证")
	}
	if err != nil {
		return false, err
	}

	ok, err := ProcessLoginRsp(cli, rsp)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func SetRelogin(cli *client.QQClient, retryInterval int, retryCount int) {
	LoginTokens.Store(cli.Uin, cli.GenToken())
	cli.DisconnectedEvent.Subscribe(func(bot *client.QQClient, e *client.ClientDisconnectedEvent) {
		if bot.Online.Load() {
			return
		}
		bot.Disconnect()
		var times = 1
		for IsClientExist(bot.Uin) {
			if bot.Online.Load() {
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

			if token, ok := LoginTokens.Load(bot.Uin); ok {
				// 尝试token登录
				if err := bot.TokenLogin(token); err != nil {
					log.Errorf("failed to relogin with token, try to login with password, %+v", err)
					bot.Disconnect()
				} else {
					LoginTokens.Store(bot.Uin, bot.GenToken())
					log.Info("succeed to relogin with token")
					return
				}
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
				LoginTokens.Store(bot.Uin, bot.GenToken())
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
	Clients.Delete(cli.Uin) // 必须先删Clients，影响IsClientExist
	LoginTokens.Delete(cli.Uin)
	if wsServers, ok := RemoteServers.Load(cli.Uin); ok {
		for _, wsServer := range wsServers {
			wsServer.Close()
		}
	}
	RemoteServers.Delete(cli.Uin)
}

func IsClientExist(uin int64) bool {
	_, ok := Clients.Load(uin)
	return ok
}

func Energy(uin uint64, id string, appVersion string, salt []byte) ([]byte, error) {
	signServer := GTL.SignServer
	if !strings.HasSuffix(signServer, "/") {
		signServer += "/"
	}
	response, err := download.Request{
		Method: http.MethodGet,
		URL:    signServer + "custom_energy" + fmt.Sprintf("?data=%v&salt=%v", id, hex.EncodeToString(salt)),
	}.Bytes()
	if err != nil {
		log.Warnf("获取T544 sign时出现错误: %v server: %v", err, signServer)
		return nil, err
	}
	data, err := hex.DecodeString(gjson.GetBytes(response, "data").String())
	if err != nil {
		log.Warnf("获取T544 sign时出现错误: %v", err)
		return nil, err
	}
	if len(data) == 0 {
		log.Warnf("获取T544 sign时出现错误: %v", "data is empty")
		return nil, errors.New("data is empty")
	}
	return data, nil
}

func Sign(seq uint64, uin string, cmd string, qua string, buff []byte) (sign []byte, extra []byte, token []byte, err error) {
	signServer := GTL.SignServer
	if !strings.HasSuffix(signServer, "/") {
		signServer += "/"
	}
	response, err := download.Request{
		Method: http.MethodPost,
		URL:    signServer + "sign",
		Header: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		Body:   bytes.NewReader([]byte(fmt.Sprintf("uin=%v&qua=%s&cmd=%s&seq=%v&buffer=%v", uin, qua, cmd, seq, hex.EncodeToString(buff)))),
	}.Bytes()
	if err != nil {
		log.Warnf("获取sso sign时出现错误: %v server: %v", err, signServer)
		return nil, nil, nil, err
	}
	sign, _ = hex.DecodeString(gjson.GetBytes(response, "data.sign").String())
	extra, _ = hex.DecodeString(gjson.GetBytes(response, "data.extra").String())
	token, _ = hex.DecodeString(gjson.GetBytes(response, "data.token").String())
	return sign, extra, token, nil
}
