package bot

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
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
	Clients                    ClientMap
	LoginTokens                TokenMap
	EnergyCount                = 0
	EnergyStop                 = false
	SignCount                  = 0
	SignStop                   = false
	RegisterSignCount          = 0
	RegisterSignStop           = false
	SubmitRequestCallbackCount = 0
	SubmitRequestCallbackStop  = false
	RequestTokenCount          = 0
	RequestTokenStop           = false
	DestoryInstanceCount       = 0
	DestoryInstanceStop        = false
	RSR                        RequestSignResult
	GTL                        *GMCLogin
	SR                         SignRegister
	IsRequestTokenAgain        bool = false
	TTI_i                           = 30
)

type Logger struct {
}

type GMCLogin struct {
	DeviceSeed     int64
	ClientProtocol int32
	SignServer     string
	SignServerKey  string
}

type SignRegister struct {
	Uin       uint64
	AndroidId string
	Guid      string
	Qimei36   string
	Key       string
}

type RequestCallback struct {
	Cmd        string `json:"cmd,omitempty"` // trpc.o3.ecdh_access.EcdhAccess.SsoSecureA2Establish
	Body       string `json:"body,omitempty"`
	CallBackId int    `json:"callbackId,omitempty"`
}

type RequestSignData struct {
	Token           string `json:"token,omitempty"`
	Extra           string `json:"extra,omitempty"`
	Sign            string `json:"sign,omitempty"`
	O3dId           string `json:"o3did,omitempty"`
	RequestCallback []*RequestCallback
}

type RequestSignResult struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Data *RequestSignData
}

func GmcTokenLogin() (g GMCLogin, err error) {
	if PathExists("deviceInfo.toml") {
		_, err = toml.DecodeFile("deviceInfo.toml", &GTL)
		return *GTL, err
	} else {
		g = GMCLogin{}
		return g, nil
	}
}

func SRI() (sr SignRegister, err error) {
	_, err = toml.DecodeFile("signRegisterInfo.toml", &SR)
	return SR, err
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
	log.Debug(format, args)
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
	if (GTL.ClientProtocol == 1 || GTL.ClientProtocol == 6) && GTL.SignServer != "" {
		RegisterSign(uint64(cli.Uin), cli.Device().AndroidId, cli.Device().Guid, cli.Device().QImei36, GTL.SignServerKey)
		if !RegisterSignStop {
			wrapper.DandelionEnergy = Energy
			wrapper.FekitGetSign = Sign
		}
	} else if GTL.SignServer != "" {
		log.Warn("SignServer 不支持该协议")
	}
	if !RegisterSignStop {
		rsp, err := cli.Login()
		if rsp.Code == byte(45) && GTL.SignServer == "" {
			log.Warn("您的账号被限制登录，请配置 SignServer 后重试")
		}
		if rsp.Code == byte(235) {
			log.Warn("设备信息被封禁，请删除设备（device）文件夹里对应设备文件后重试")
		}
		if rsp.Code == byte(237) {
			log.Warn("登录过于频繁，请在手机QQ登录并根据提示完成认证")
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
	return false, fmt.Errorf("登录失败！")
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
	DestoryInstance(uint(SR.Uin), SR.Key)
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
	if !EnergyStop {
		signServer := GTL.SignServer
		if !strings.HasSuffix(signServer, "/") {
			signServer += "/"
		}
		response, err := download.Request{
			Method: http.MethodGet,
			URL:    signServer + "custom_energy" + fmt.Sprintf("?uin=%v&data=%v&salt=%v", uin, id, hex.EncodeToString(salt)),
		}.Bytes()
		if err != nil {
			log.Warnf("获取T544 sign时出现错误: %v server: %v", err, signServer)
			EnergyCount++
			if EnergyCount > 2 {
				EnergyStop = true
			}
			return nil, err
		}
		data, err := hex.DecodeString(gjson.GetBytes(response, "data").String())
		if err != nil {
			log.Warnf("获取T544 sign时出现错误: %v", err)
			EnergyCount++
			if EnergyCount > 2 {
				EnergyStop = true
			}
			return nil, err
		}
		if len(data) == 0 {
			log.Warnf("获取T544 sign时出现错误: %v", "data is empty")
			EnergyCount++
			if EnergyCount > 2 {
				EnergyStop = true
			}
			return nil, errors.New("data is empty")
		}
		return data, nil
	} else {
		log.Warn("Energy失败，请重试")
		DestoryInstance(uint(uin), SR.Key)
		return nil, fmt.Errorf("Energy失败，请重试")
	}
}

func Sign(seq uint64, uin string, cmd string, qua string, buff []byte) (sign []byte, extra []byte, token []byte, err error) {
	if !SignStop {
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
			SignCount++
			if SignCount > 2 {
				SignStop = true
			}
			return nil, nil, nil, err
		}
		sign, _ = hex.DecodeString(gjson.GetBytes(response, "data.sign").String())
		extra, _ = hex.DecodeString(gjson.GetBytes(response, "data.extra").String())
		token, _ = hex.DecodeString(gjson.GetBytes(response, "data.token").String())

		json.Unmarshal(response, &RSR)
		if len(RSR.Data.RequestCallback) > 1 {
			log.Warn(RSR.Data.RequestCallback[0], RSR.Data.RequestCallback[1])
		}
		return sign, extra, token, nil
	} else {
		log.Warn("Sign失败, 请重试")
		DestoryInstance(uint(SR.Uin), SR.Key)
		return nil, nil, nil, fmt.Errorf("Sign失败, 请重试")
	}
}

func RegisterSign(uin uint64, androidId []byte, guid []byte, Qimei36 string, signServerAuth string) {
	signServer := GTL.SignServer
	if !strings.HasSuffix(signServer, "/") {
		signServer += "/"
	}
	SR.Uin = uin
	SR.AndroidId = string(androidId)
	SR.Guid = string(guid)
	SR.Guid = string(guid)
	SR.Key = signServerAuth
	// http://your.host:port/register?uin=[QQ]&android_id=[ANDROID_ID]&guid=[GUID]&qimei36=[QIMEI36]&key=[KEY]
	_ = os.WriteFile("signRegisterInfo.toml", []byte(fmt.Sprintf("uin= %v \nandroidId= \"%s\" \nguid= \"%s\" \nqimei36= \"%s\" \nkey= \"%s\"", uin, hex.EncodeToString(androidId), hex.EncodeToString(guid), Qimei36, signServerAuth)), 0o644)

	log.Warn(uin, hex.EncodeToString(androidId), hex.EncodeToString(guid), Qimei36, signServerAuth)
	log.Warn(fmt.Sprintf("?uin=%v&android_id=%s&guid=%s&qimei36=%s&key=%s", uin, hex.EncodeToString(androidId), hex.EncodeToString(guid), Qimei36, signServerAuth))
	response, err := download.Request{
		Method: http.MethodGet,
		URL:    signServer + "register" + fmt.Sprintf("?uin=%v&android_id=%s&guid=%s&qimei36=%s&key=%s", uin, hex.EncodeToString(androidId), hex.EncodeToString(guid), Qimei36, signServerAuth),
	}.Bytes()
	if err != nil {
		log.Warnf("初始化 Sign 失败\n", err)
		if RegisterSignCount < 2 {
			time.Sleep(time.Second * 5)
			RegisterSign(SR.Uin, []byte(SR.AndroidId), []byte(SR.Guid), SR.Qimei36, signServerAuth)
			RegisterSignCount++
		} else {
			RegisterSignStop = true
		}
	} else {
		log.Info("初始化 Sign 成功")
		log.Warn(gjson.GetBytes(response, "msg").String())
	}
}

// http://your.host:port/submit?uin=[QQ]&cmd=[CMD]&callback_id=[CALLBACK_ID]&buffer=[BUFFER]
func SubmitRequestCallback(uin uint64, cmd string, callbackId int, buffer []byte) {
	if !SubmitRequestCallbackStop {
		signServer := GTL.SignServer
		if !strings.HasSuffix(signServer, "/") {
			signServer += "/"
		}
		response, err := download.Request{
			Method: http.MethodGet,
			URL:    signServer + "submit" + fmt.Sprintf("?uin=%v&cmd=%s&callback_id=%v&buffer=%s", uin, cmd, callbackId, buffer),
		}.Bytes()
		if err != nil {
			log.Warnf(cmd, " ", callbackId, "提交失败\n", err)
			if SubmitRequestCallbackCount < 2 {
				time.Sleep(time.Second * 5)
				SubmitRequestCallback(uin, cmd, callbackId, buffer)
				SubmitRequestCallbackCount++
			} else {
				SubmitRequestCallbackStop = true
			}
		} else {
			log.Info(cmd, " ", callbackId, "提交成功")
			log.Warn(string(response))
			log.Warn(gjson.GetBytes(response, "msg").String())
		}
	} else {
		log.Warn("SubmitRequestCallback失败，请重试")
		DestoryInstance(uint(SR.Uin), SR.Key)
	}
}

func RequestToken(uin uint64) {
	if !RequestTokenStop {
		signServer := GTL.SignServer
		if !strings.HasSuffix(signServer, "/") {
			signServer += "/"
		}
		response, err := download.Request{
			Method: http.MethodGet,
			URL:    signServer + "request_token" + fmt.Sprintf("?uin=%v", uin),
		}.Bytes()
		if err != nil || strings.HasPrefix(gjson.GetBytes(response, "msg").String(), "Uin") { // QSign
			log.Warnf("请求 Token 失败\n", gjson.GetBytes(response, "msg").String(), err)
			log.Info("正在重新注册 ", uin)
			if RequestTokenCount < 2 {
				time.Sleep(time.Second * 5)
				RegisterSign(SR.Uin, []byte(SR.AndroidId), []byte(SR.Guid), SR.Qimei36, SR.Key)
				IsRequestTokenAgain = true
				RequestTokenCount++
			} else {
				RequestTokenStop = true
			}
		} else if strings.HasPrefix(gjson.GetBytes(response, "msg").String(), "QSign") {
			log.Warn("QSign not initialized, unable to request_ Token, please submit the initialization package first.")
			IsRequestTokenAgain = false
		} else {
			log.Info("请求 Token 成功")
			log.Warn(string(response))
			log.Warn(gjson.GetBytes(response, "msg").String())
		}
	} else {
		log.Warn("RequestToken失败，请重试！")
		DestoryInstance(uint(SR.Uin), SR.Key)
	}
}

// http://host:port/destroy?uin=[QQ]&key=[key]
func DestoryInstance(uin uint, key string) {
	signServer := GTL.SignServer
	if !strings.HasSuffix(signServer, "/") {
		signServer += "/"
	}
	response, err := download.Request{
		Method: http.MethodGet,
		URL:    signServer + "destroy" + fmt.Sprintf("?uin=%v&key=%s", uin, key),
	}.Bytes()
	if err != nil {
		if DestoryInstanceCount < 2 {
			time.Sleep(time.Second * 5)
			DestoryInstance(uin, key)
			DestoryInstanceCount++
		} else {
			DestoryInstanceStop = true
		}
	} else {
		log.Warn(gjson.GetBytes(response, "msg").String())
	}
}

func TTIR(uin uint64) {
	for TTI_i >= 0 {
		if RequestTokenStop {
			break
		}
		time.Sleep(time.Minute)
		if TTI_i == 0 {
			TTI_i = 30
			RequestToken(uin)
		}
		TTI_i--
	}
}
