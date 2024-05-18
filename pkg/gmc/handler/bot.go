package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/2mf8/Go-Lagrange-Client/pkg/bot"
	"github.com/2mf8/Go-Lagrange-Client/pkg/config"
	"github.com/2mf8/Go-Lagrange-Client/pkg/device"
	"github.com/2mf8/Go-Lagrange-Client/pkg/gmc/plugins"
	"github.com/2mf8/Go-Lagrange-Client/pkg/plugin"
	"github.com/2mf8/Go-Lagrange-Client/pkg/util"
	"github.com/2mf8/Go-Lagrange-Client/proto_gen/dto"

	"github.com/2mf8/LagrangeGo/client"
	"github.com/2mf8/LagrangeGo/client/auth"
	_ "github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

var queryQRCodeMutex = &sync.RWMutex{}
var qrCodeBot *client.QQClient

type QRCodeResp int

const (
	Unknown = iota
	QRCodeImageFetch
	QRCodeWaitingForScan
	QRCodeWaitingForConfirm
	QRCodeTimeout
	QRCodeConfirmed
	QRCodeCanceled
)

func TokenLogin() {
	dfs, err := os.ReadDir("./device/")
	if err == nil {
		for _, v := range dfs {
			df := strings.Split(v.Name(), ".")
			uin, err := strconv.ParseInt(df[0], 10, 64)
			if err == nil {
				devi := device.GetDevice(uin)
				sfs, err := os.ReadDir("./sig/")
				if err == nil {
					for _, sv := range sfs {
						sf := strings.Split(sv.Name(), ".")
						if df[0] == sf[0] {
							sigpath := fmt.Sprintf("./sig/%s", sv.Name())
							data, err := os.ReadFile(sigpath)
							if err == nil {
								sig, err := auth.UnmarshalSigInfo(data, true)
								if err == nil {
									go func() {
										queryQRCodeMutex.Lock()
										defer queryQRCodeMutex.Unlock()
										appInfo := auth.AppList["linux"]
										qrCodeBot = client.NewClient(0, "https://sign.lagrangecore.org/api/sign", appInfo)
										qrCodeBot.UseDevice(devi)
										qrCodeBot.UseSig(sig)
										qrCodeBot.SessionLogin()
										go AfterLogin(qrCodeBot)
									}()
								}
							}
						}
					}
				}
			} else {
				fmt.Printf("转换账号%s失败", df[0])
			}
		}
	}
}

func init() {
	log.Infof("加载日志插件 Log")
	plugin.AddPrivateMessagePlugin(plugins.LogPrivateMessage)
	plugin.AddGroupMessagePlugin(plugins.LogGroupMessage)
	log.Infof("加载测试插件 Hello")
	plugin.AddPrivateMessagePlugin(plugins.HelloPrivateMessage)
	plugin.AddGroupMessagePlugin(plugins.HelloGroupMessage)
	log.Infof("加载上报插件 Report")
	plugin.AddPrivateMessagePlugin(plugins.ReportPrivateMessage)
	plugin.AddGroupMessagePlugin(plugins.ReportGroupMessage)
	plugin.AddMemberJoinGroupPlugin(plugins.ReportMemberJoin)
	plugin.AddMemberLeaveGroupPlugin(plugins.ReportMemberLeave)
	plugin.AddNewFriendRequestPlugin(plugins.ReportNewFriendRequest)
	plugin.AddGroupInvitedRequestPlugin(plugins.ReportGroupInvitedRequest)
	plugin.AddGroupMessageRecalledPlugin(plugins.ReportGroupMessageRecalled)
	plugin.AddFriendMessageRecalledPlugin(plugins.ReportFriendMessageRecalled)
	plugin.AddNewFriendAddedPlugin(plugins.ReportNewFriendAdded)
	plugin.AddGroupMutePlugin(plugins.ReportGroupMute)
}

func DeleteBot(c *gin.Context) {
	sigpath := fmt.Sprintf("./sig/%v.bin", qrCodeBot.Uin)
	sigDir := path.Dir(sigpath)
	if !util.PathExists(sigDir) {
		log.Infof("%+v 目录不存在，自动创建", sigDir)
		if err := os.MkdirAll(sigDir, 0777); err != nil {
			log.Warnf("failed to mkdir deviceDir, err: %+v", err)
		}
	}
	data, err := qrCodeBot.Sig().Marshal()
	if err != nil {
		log.Errorln("marshal sig.bin err:", err)
		return
	}
	err = os.WriteFile(sigpath, data, 0644)
	if err != nil {
		log.Errorln("write sig.bin err:", err)
		return
	}
	log.Infoln("sig saved into sig.bin")
	req := &dto.DeleteBotReq{}
	err = Bind(c, req)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request, not protobuf")
		return
	}
	cli, ok := bot.Clients.Load(req.BotId)
	if !ok {
		c.String(http.StatusBadRequest, "bot not exists")
		return
	}
	bot.Clients.Delete(int64(cli.Uin))
	bot.ReleaseClient(cli)
	resp := &dto.DeleteBotResp{}
	Return(c, resp)
}

func ListBot(c *gin.Context) {
	req := &dto.ListBotReq{}
	err := Bind(c, req)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request, not protobuf")
		return
	}
	var resp = &dto.ListBotResp{
		BotList: []*dto.Bot{},
	}
	bot.Clients.Range(func(_ int64, cli *client.QQClient) bool {
		resp.BotList = append(resp.BotList, &dto.Bot{
			BotId:    int64(cli.Uin),
			IsOnline: cli.Online.Load(),
		})
		return true
	})
	Return(c, resp)
}

func FetchQrCode(c *gin.Context) {
	req := &dto.FetchQRCodeReq{}
	err := Bind(c, req)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request, not protobuf")
		return
	}
	newDeviceInfo := device.GetDevice(req.DeviceSeed)
	appInfo := auth.AppList["linux"]
	if err != nil {
		fmt.Println(err)
	} else {
		qqclient := client.NewClient(0, "https://sign.lagrangecore.org/api/sign", appInfo)
		qqclient.UseDevice(newDeviceInfo)
		qrCodeBot = qqclient
		b, s, err := qrCodeBot.FecthQRCode()
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("failed to fetch qrcode, %+v", err))
			return
		}
		resp := &dto.QRCodeLoginResp{
			State:     dto.QRCodeLoginResp_QRCodeLoginState(http.StatusOK),
			ImageData: b,
			Sig:       []byte(s),
		}
		Return(c, resp)
	}
}

func QueryQRCodeStatus(c *gin.Context) {
	queryQRCodeMutex.Lock()
	defer queryQRCodeMutex.Unlock()
	respCode := 0
	ok, err := qrCodeBot.GetQRCodeResult()
	if err != nil {
		resp := &dto.QRCodeLoginResp{
			State: dto.QRCodeLoginResp_QRCodeLoginState(http.StatusExpectationFailed),
		}
		Return(c, resp)
	}
	fmt.Println(ok.Name())
	if !ok.Success() {
		resp := &dto.QRCodeLoginResp{
			State: dto.QRCodeLoginResp_QRCodeLoginState(http.StatusExpectationFailed),
		}
		Return(c, resp)
	}
	if ok.Name() == "WaitingForConfirm" {
		respCode = QRCodeWaitingForScan
	}
	if ok.Name() == "Canceled" {
		respCode = QRCodeCanceled
	}
	if ok.Name() == "WaitingForConfirm" {
		respCode = QRCodeWaitingForConfirm
	}
	if ok.Name() == "Confirmed" {
		respCode = QRCodeConfirmed
		err := qrCodeBot.QRCodeConfirmed()
		if err == nil {
			go func() {
				queryQRCodeMutex.Lock()
				defer queryQRCodeMutex.Unlock()
				qrCodeBot.Init()
				time.Sleep(time.Second * 5)
				AfterLogin(qrCodeBot)
			}()
		}
	}
	if ok.Name() == "Expired" {
		respCode = QRCodeTimeout
	}
	resp := &dto.QRCodeLoginResp{
		State: dto.QRCodeLoginResp_QRCodeLoginState(respCode),
	}
	Return(c, resp)
}

func ListPlugin(c *gin.Context) {
	req := &dto.ListPluginReq{}
	err := Bind(c, req)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	var resp = &dto.ListPluginResp{
		Plugins: []*dto.Plugin{},
	}
	config.Plugins.Range(func(key string, p *config.Plugin) bool {
		resp.Plugins = append(resp.Plugins, &dto.Plugin{
			Name:         p.Name,
			Disabled:     p.Disabled,
			Json:         p.Json,
			Protocol: 	  p.Protocol,
			Urls:         p.Urls,
			EventFilter:  p.EventFilter,
			ApiFilter:    p.ApiFilter,
			RegexFilter:  p.RegexFilter,
			RegexReplace: p.RegexReplace,
			ExtraHeader: func() []*dto.Plugin_Header {
				headers := make([]*dto.Plugin_Header, 0)
				for k, v := range p.ExtraHeader {
					headers = append(headers, &dto.Plugin_Header{
						Key:   k,
						Value: v,
					})
				}
				return headers
			}(),
		})
		return true
	})
	Return(c, resp)
}

func SavePlugin(c *gin.Context) {
	req := &dto.SavePluginReq{}
	err := Bind(c, req)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	if req.Plugin == nil {
		c.String(http.StatusBadRequest, "plugin is nil")
		return
	}
	p := req.Plugin
	if p.ApiFilter == nil {
		p.ApiFilter = []int32{}
	}
	if p.EventFilter == nil {
		p.EventFilter = []int32{}
	}
	if p.Urls == nil {
		p.Urls = []string{}
	}
	config.Plugins.Store(p.Name, &config.Plugin{
		Name:         p.Name,
		Disabled:     p.Disabled,
		Json:         p.Json,
		Protocol:     p.Protocol,
		Urls:         p.Urls,
		EventFilter:  p.EventFilter,
		ApiFilter:    p.ApiFilter,
		RegexFilter:  p.RegexFilter,
		RegexReplace: p.RegexReplace,
		ExtraHeader: func() map[string][]string {
			headers := map[string][]string{}
			for _, h := range p.ExtraHeader {
				headers[h.Key] = h.Value
			}
			return headers
		}(),
	})
	config.WritePlugins()
	resp := &dto.SavePluginResp{}
	Return(c, resp)
}

func DeletePlugin(c *gin.Context) {
	req := &dto.DeletePluginReq{}
	err := Bind(c, req)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	config.Plugins.Delete(req.Name)
	config.WritePlugins()
	resp := &dto.DeletePluginResp{}
	Return(c, resp)
}

func Return(c *gin.Context, resp proto.Message) {
	var (
		data []byte
		err  error
	)
	switch c.ContentType() {
	case binding.MIMEPROTOBUF:
		data, err = proto.Marshal(resp)
	case binding.MIMEJSON:
		data, err = json.Marshal(resp)
	}
	if err != nil {
		c.String(http.StatusInternalServerError, "marshal resp error")
		return
	}
	c.Data(http.StatusOK, c.ContentType(), data)
}

func AfterLogin(cli *client.QQClient) {
	for {
		time.Sleep(5 * time.Second)
		if cli.Online.Load() {
			break
		}
		log.Warnf("机器人不在线，可能在等待输入验证码，或出错了。如果出错请重启。")
	}
	bot.Clients.Store(int64(cli.Uin), cli)
	plugin.Serve(cli)
	log.Infof("插件加载完成")

	log.Infof("刷新好友列表")
	if fs, err := cli.GetFriendsData(); err != nil {
		util.FatalError(fmt.Errorf("failed to load friend list, err: %+v", err))
	} else {
		log.Infof("共加载 %v 个好友.", len(fs))
	}

	bot.ConnectUniversal(cli)

	defer cli.Release()
	defer func() {
		sigpath := fmt.Sprintf("./sig/%v.bin", cli.Uin)
		sigDir := path.Dir(sigpath)
		if !util.PathExists(sigDir) {
			log.Infof("%+v 目录不存在，自动创建", sigDir)
			if err := os.MkdirAll(sigDir, 0777); err != nil {
				log.Warnf("failed to mkdir deviceDir, err: %+v", err)
			}
		}
		data, err := qrCodeBot.Sig().Marshal()
		if err != nil {
			log.Errorln("marshal sig.bin err:", err)
			return
		}
		err = os.WriteFile(sigpath, data, 0644)
		if err != nil {
			log.Errorln("write sig.bin err:", err)
			return
		}
		log.Infoln("sig saved into sig.bin")
	}()

	// setup the main stop channel
	mc := make(chan os.Signal, 2)
	signal.Notify(mc, os.Interrupt, syscall.SIGTERM)
	for {
		switch <-mc {
		case os.Interrupt, syscall.SIGTERM:
			return
		}
	}
}
