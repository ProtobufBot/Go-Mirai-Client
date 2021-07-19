package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ProtobufBot/Go-Mirai-Client/pkg/device"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/plugin"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/util"
	"github.com/ProtobufBot/Go-Mirai-Client/proto_gen/dto"
	"github.com/ProtobufBot/Go-Mirai-Client/service/bot"
	"github.com/ProtobufBot/Go-Mirai-Client/service/plugins"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

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
	_, ok := bot.Clients[req.BotId]
	if ok {
		c.String(http.StatusInternalServerError, "botId already exists")
	}
	go func() {
		CreateBotImpl(req.BotId, req.Password, req.DeviceSeed)
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
	cli, ok := bot.Clients[req.BotId]
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
	for _, cli := range bot.Clients {
		resp.BotList = append(resp.BotList, &dto.Bot{
			BotId:    cli.Uin,
			IsOnline: cli.Online,
		})
	}
	Return(c, resp)
}

func ListCaptcha(c *gin.Context) {
	req := &dto.ListCaptchaReq{}
	err := c.Bind(req)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request, not protobuf")
		return
	}
	var resp *dto.ListCaptchaResp
	if bot.Captcha != nil {
		resp = &dto.ListCaptchaResp{
			CaptchaList: []*dto.Captcha{bot.Captcha},
		}
	} else {
		resp = &dto.ListCaptchaResp{
			CaptchaList: []*dto.Captcha{},
		}
	}
	Return(c, resp)
}

func SolveCaptcha(c *gin.Context) {
	req := &dto.SolveCaptchaReq{}
	err := c.Bind(req)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request, not protobuf")
		return
	}
	captchaProm, ok := bot.CaptchaPromises[req.BotId]
	if !ok {
		c.String(http.StatusInternalServerError, "captcha not found")
		return
	}

	err = captchaProm.Resolve(req.Result)
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
	cli, ok := bot.Clients[0]
	if ok {
		cli.Release()
	}
	cli = client.NewClientEmpty()
	deviceInfo := device.GetDevice(req.DeviceSeed)
	cli.UseDevice(deviceInfo)
	bot.Clients[0] = cli

	log.Infof("初始化日志")
	bot.InitLog(cli)
	fetchQRCodeResp, err := cli.FetchQRCode()
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
	req := &dto.QueryQRCodeStatusReq{}
	err := c.Bind(req)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("failed to bind, %+v", err))
		return
	}

	cli, ok := bot.Clients[req.BotId]
	if !ok {
		c.String(http.StatusBadRequest, "botId not exists")
		return
	}

	if cli.Online {
		c.String(http.StatusBadRequest, "already online")
		return
	}

	queryQRCodeStatusResp, err := cli.QueryQRCodeStatus(req.Sig)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("failed to query qrcode status, %+v", err))
		return
	}
	if queryQRCodeStatusResp.State == client.QRCodeConfirmed {
		loginResp, err := cli.QRCodeLogin(queryQRCodeStatusResp.LoginInfo)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("failed to qrcode login, %+v", err))
			return
		}
		go func() {
			ok, err := bot.ProcessLoginRsp(cli, loginResp)
			if err != nil {
				util.FatalError(fmt.Errorf("failed to login, err: %+v", err))
			}
			if ok {
				log.Infof("登录成功")
				originCli, ok := bot.Clients[cli.Uin]
				if ok {
					// 重复登录，旧的断开
					originCli.Release()
				}
				delete(bot.Clients, 0)
				bot.Clients[cli.Uin] = cli
				AfterLogin(cli)
			} else {
				log.Infof("登录失败")
			}
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

func CreateBotImpl(uin int64, password string, deviceRandSeed int64) {
	log.Infof("开始初始化设备信息")
	deviceInfo := device.GetDevice(uin)
	if deviceRandSeed != 0 {
		deviceInfo = device.GetDevice(deviceRandSeed)
	}
	log.Infof("设备信息 %+v", string(deviceInfo.ToJson()))

	log.Infof("创建机器人 %+v", uin)

	cli := client.NewClient(uin, password)
	cli.UseDevice(deviceInfo)
	bot.Clients[uin] = cli

	log.Infof("初始化日志")
	bot.InitLog(cli)

	log.Infof("登录中...")
	ok, err := bot.Login(cli)
	if err != nil {
		util.FatalError(fmt.Errorf("failed to login, err: %+v", err))
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
		if cli.Online {
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
