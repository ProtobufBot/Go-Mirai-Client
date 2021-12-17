package handler

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ProtobufBot/Go-Mirai-Client/pkg/bot"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/device"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/gmc/plugins"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/plugin"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/util"
	"github.com/ProtobufBot/Go-Mirai-Client/proto_gen/dto"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

var queryQRCodeMutex = &sync.RWMutex{}
var qrCodeBot *client.QQClient

func init() {
	//log.Infof("加载日志插件 Log")
	plugin.AddPrivateMessagePlugin(plugins.LogPrivateMessage)
	plugin.AddGroupMessagePlugin(plugins.LogGroupMessage)

	//log.Infof("加载测试插件 Hello")
	plugin.AddPrivateMessagePlugin(plugins.HelloPrivateMessage)

	//log.Infof("加载上报插件 Report")
	plugin.AddPrivateMessagePlugin(plugins.ReportPrivateMessage)
	plugin.AddGroupMessagePlugin(plugins.ReportGroupMessage)
	plugin.AddTempMessagePlugin(plugins.ReportTempMessage)
	plugin.AddMemberPermissionChangedPlugin(plugins.ReportMemberPermissionChanged)
	plugin.AddMemberJoinGroupPlugin(plugins.ReportMemberJoin)
	plugin.AddMemberLeaveGroupPlugin(plugins.ReportMemberLeave)
	plugin.AddJoinGroupPlugin(plugins.ReportJoinGroup)
	plugin.AddLeaveGroupPlugin(plugins.ReportLeaveGroup)
	plugin.AddNewFriendRequestPlugin(plugins.ReportNewFriendRequest)
	plugin.AddUserJoinGroupRequestPlugin(plugins.ReportUserJoinGroupRequest)
	plugin.AddGroupInvitedRequestPlugin(plugins.ReportGroupInvitedRequest)
	plugin.AddGroupMessageRecalledPlugin(plugins.ReportGroupMessageRecalled)
	plugin.AddFriendMessageRecalledPlugin(plugins.ReportFriendMessageRecalled)
	plugin.AddNewFriendAddedPlugin(plugins.ReportNewFriendAdded)
	plugin.AddOfflineFilePlugin(plugins.ReportOfflineFile)
	plugin.AddGroupMutePlugin(plugins.ReportGroupMute)
}

func CreateBot(c *gin.Context) {
	req := &dto.CreateBotReq{}
	err := c.Bind(req)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request, not protobuf")
		return
	}
	if req.BotId == 0 {
		c.String(http.StatusBadRequest, "botId is 0")
		return
	}
	_, ok := bot.Clients.Load(req.BotId)
	if ok {
		c.String(http.StatusInternalServerError, "botId already exists")
		return
	}
	go func() {
		CreateBotImpl(req.BotId, req.Password, req.DeviceSeed, req.ClientProtocol)
	}()
	resp := &dto.CreateBotResp{}
	Return(c, resp)
}

func DeleteBot(c *gin.Context) {
	req := &dto.DeleteBotReq{}
	err := c.Bind(req)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request, not protobuf")
		return
	}
	cli, ok := bot.Clients.Load(req.BotId)
	if !ok {
		c.String(http.StatusBadRequest, "bot not exists")
		return
	}
	bot.ReleaseClient(cli)
	resp := &dto.DeleteBotResp{}
	Return(c, resp)
}

func ListBot(c *gin.Context) {
	req := &dto.ListBotReq{}
	err := c.Bind(req)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request, not protobuf")
		return
	}
	var resp = &dto.ListBotResp{
		BotList: []*dto.Bot{},
	}
	bot.Clients.Range(func(_ int64, cli *client.QQClient) bool {
		resp.BotList = append(resp.BotList, &dto.Bot{
			BotId:    cli.Uin,
			IsOnline: cli.Online.Load(),
			Captcha: func() *dto.Bot_Captcha {
				if waitingCaptcha, ok := bot.WaitingCaptchas.Load(cli.Uin); ok {
					return waitingCaptcha.Captcha
				}
				return nil
			}(),
		})
		return true
	})
	Return(c, resp)
}

func SolveCaptcha(c *gin.Context) {
	req := &dto.SolveCaptchaReq{}
	err := c.Bind(req)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request, not protobuf")
		return
	}
	waitingCaptcha, ok := bot.WaitingCaptchas.Load(req.BotId)
	if !ok {
		c.String(http.StatusInternalServerError, "captcha not found")
		return
	}

	err = waitingCaptcha.Prom.Resolve(req.Result)
	if err != nil {
		c.String(http.StatusInternalServerError, "solve captcha error")
		return
	}

	resp := &dto.SolveCaptchaResp{}
	Return(c, resp)
}

func FetchQrCode(c *gin.Context) {
	req := &dto.FetchQRCodeReq{}
	err := c.Bind(req)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request, not protobuf")
		return
	}
	if qrCodeBot != nil {
		qrCodeBot.Release()
	}
	qrCodeBot = client.NewClientEmpty()
	deviceInfo := device.GetDevice(req.DeviceSeed, req.ClientProtocol)
	qrCodeBot.UseDevice(deviceInfo)

	log.Infof("初始化日志")
	bot.InitLog(qrCodeBot)
	fetchQRCodeResp, err := qrCodeBot.FetchQRCode()
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("failed to fetch qrcode, %+v", err))
		return
	}
	resp := &dto.QRCodeLoginResp{
		State:     dto.QRCodeLoginResp_QRCodeLoginState(fetchQRCodeResp.State),
		ImageData: fetchQRCodeResp.ImageData,
		Sig:       fetchQRCodeResp.Sig,
	}
	Return(c, resp)
}

func QueryQRCodeStatus(c *gin.Context) {
	queryQRCodeMutex.Lock()
	defer queryQRCodeMutex.Unlock()
	req := &dto.QueryQRCodeStatusReq{}
	err := c.Bind(req)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("failed to bind, %+v", err))
		return
	}

	if qrCodeBot == nil {
		c.String(http.StatusBadRequest, "please fetch qrcode first")
		return
	}

	if qrCodeBot.Online.Load() {
		c.String(http.StatusBadRequest, "already online")
		return
	}

	queryQRCodeStatusResp, err := qrCodeBot.QueryQRCodeStatus(req.Sig)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("failed to query qrcode status, %+v", err))
		return
	}
	if queryQRCodeStatusResp.State == client.QRCodeConfirmed {
		go func() {
			queryQRCodeMutex.Lock()
			defer queryQRCodeMutex.Unlock()

			loginResp, err := qrCodeBot.QRCodeLogin(queryQRCodeStatusResp.LoginInfo)
			if err != nil {
				c.String(http.StatusInternalServerError, fmt.Sprintf("failed to qrcode login, %+v", err))
				return
			}
			if !loginResp.Success {
				c.String(http.StatusInternalServerError, fmt.Sprintf("failed to qrcode login, %+v", err))
				return
			}
			log.Infof("登录成功")
			originCli, ok := bot.Clients.Load(qrCodeBot.Uin)
			// 重复登录，旧的断开
			if ok {
				originCli.Release()
			}
			bot.Clients.Store(qrCodeBot.Uin, qrCodeBot)
			go AfterLogin(qrCodeBot)
			qrCodeBot = nil
		}()
	}

	resp := &dto.QRCodeLoginResp{
		State:     dto.QRCodeLoginResp_QRCodeLoginState(queryQRCodeStatusResp.State),
		ImageData: queryQRCodeStatusResp.ImageData,
		Sig:       queryQRCodeStatusResp.Sig,
	}
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

func CreateBotImpl(uin int64, password string, deviceRandSeed int64, clientProtocol int32) {
	CreateBotImplMd5(uin, md5.Sum([]byte(password)), deviceRandSeed, clientProtocol)
}

func CreateBotImplMd5(uin int64, passwordMd5 [16]byte, deviceRandSeed int64, clientProtocol int32) {
	log.Infof("开始初始化设备信息")
	deviceInfo := device.GetDevice(uin, clientProtocol)
	if deviceRandSeed != 0 {
		deviceInfo = device.GetDevice(deviceRandSeed, clientProtocol)
	}
	log.Infof("设备信息 %+v", string(deviceInfo.ToJson()))

	log.Infof("创建机器人 %+v", uin)

	cli := client.NewClientMd5(uin, passwordMd5)
	cli.UseDevice(deviceInfo)
	bot.Clients.Store(uin, cli)

	log.Infof("初始化日志")
	bot.InitLog(cli)

	log.Infof("登录中...")
	ok, err := bot.Login(cli)
	if err != nil {
		// TODO 登录失败，是否需要删除？
		log.Errorf("failed to login, err: %+v", err)
		return
	}
	if ok {
		log.Infof("登录成功")
		AfterLogin(cli)
	} else {
		log.Infof("登录失败")
	}
}

func AfterLogin(cli *client.QQClient) {
	for {
		time.Sleep(5 * time.Second)
		if cli.Online.Load() {
			break
		}
		log.Warnf("机器人不在线，可能在等待输入验证码，或出错了。如果出错请重启。")
	}
	plugin.Serve(cli)
	log.Infof("插件加载完成")

	log.Infof("刷新好友列表")
	if err := cli.ReloadFriendList(); err != nil {
		util.FatalError(fmt.Errorf("failed to load friend list, err: %+v", err))
	}
	log.Infof("共加载 %v 个好友.", len(cli.FriendList))

	log.Infof("刷新群列表")
	if err := cli.ReloadGroupList(); err != nil {
		util.FatalError(fmt.Errorf("failed to load group list, err: %+v", err))
	}
	log.Infof("共加载 %v 个群.", len(cli.GroupList))

	bot.ConnectUniversal(cli)

	bot.SetRelogin(cli, 30, 20)
}
