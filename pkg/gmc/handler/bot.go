package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/2mf8/Go-Lagrange-Client/pkg/bot"
	"github.com/2mf8/Go-Lagrange-Client/pkg/config"
	"github.com/2mf8/Go-Lagrange-Client/pkg/gmc/plugins"
	"github.com/2mf8/Go-Lagrange-Client/pkg/plugin"
	"github.com/2mf8/Go-Lagrange-Client/pkg/util"
	"github.com/2mf8/Go-Lagrange-Client/proto_gen/dto"

	_ "github.com/BurntSushi/toml"
	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

var queryQRCodeMutex = &sync.RWMutex{}
var qrCodeBot *client.QQClient

func init() {
	log.Infof("加载日志插件 Log")
	plugin.AddPrivateMessagePlugin(plugins.LogPrivateMessage)
	plugin.AddGroupMessagePlugin(plugins.LogGroupMessage)
	log.Infof("加载测试插件 Hello")
	plugin.AddPrivateMessagePlugin(plugins.HelloPrivateMessage)
	log.Infof("加载上报插件 Report")
	plugin.AddPrivateMessagePlugin(plugins.ReportPrivateMessage)
	plugin.AddGroupMessagePlugin(plugins.ReportGroupMessage)
	plugin.AddMemberJoinGroupPlugin(plugins.ReportMemberJoin)
}

func DeleteBot(c *gin.Context) {
	req := &dto.DeleteBotReq{}
	err := Bind(c, req)
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
	b, s, err := qrCodeBot.FecthQrcode()
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

func QueryQRCodeStatus(c *gin.Context) {
	r, err := qrCodeBot.GetQrcodeResult()

	if err != nil {
		resp := &dto.QRCodeLoginResp{
			State: dto.QRCodeLoginResp_QRCodeLoginState(http.StatusExpectationFailed),
		}
		Return(c, resp)
	}

	if !r.Success() {
		resp := &dto.QRCodeLoginResp{
			State: dto.QRCodeLoginResp_QRCodeLoginState(http.StatusExpectationFailed),
		}
		Return(c, resp)
	}

	resp := &dto.QRCodeLoginResp{
		State: dto.QRCodeLoginResp_QRCodeLoginState(http.StatusOK),
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
}
