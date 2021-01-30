package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/plugin"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/util"
	"github.com/ProtobufBot/Go-Mirai-Client/proto_gen/dto"
	"github.com/ProtobufBot/Go-Mirai-Client/service/bot"
	"github.com/ProtobufBot/Go-Mirai-Client/service/plugins"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

func CreateBot(c *gin.Context) {
	req := &dto.CreateBotReq{}
	err := c.Bind(req)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request, not protobuf")
		return
	}
	if bot.Cli != nil {
		c.String(http.StatusInternalServerError, "only one bot is allowed")
	}
	go func() {
		CreateBotImpl(req.BotId, req.Password)
	}()
	resp := &dto.CreateBotResp{}
	Return(c, resp)
}

func ListBot(c *gin.Context) {
	req := &dto.ListBotReq{}
	err := c.Bind(req)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request, not protobuf")
		return
	}
	var resp *dto.ListBotResp
	if bot.Cli != nil {
		resp = &dto.ListBotResp{
			BotList: []*dto.Bot{
				{
					BotId:    bot.Cli.Uin,
					IsOnline: bot.Cli.Online,
				},
			},
		}
	} else {
		resp = &dto.ListBotResp{
			BotList: []*dto.Bot{},
		}
	}
	Return(c, resp)
}

func ListCaptcha(c *gin.Context) {
	req := &dto.GetCaptchaListReq{}
	err := c.Bind(req)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request, not protobuf")
		return
	}
	var resp *dto.GetCaptchaListResp
	if bot.Captcha != nil {
		resp = &dto.GetCaptchaListResp{
			CaptchaList: []*dto.Captcha{bot.Captcha},
		}
	} else {
		resp = &dto.GetCaptchaListResp{
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
	if bot.CaptchaPromise == nil {
		c.String(http.StatusInternalServerError, "captcha not found")
		return
	}

	err = bot.CaptchaPromise.Resolve(req.Result)
	if err != nil {
		c.String(http.StatusInternalServerError, "solve captcha error")
		return
	}

	resp := &dto.SolveCaptchaResp{}
	Return(c, resp)
}

func Return(c *gin.Context, resp proto.Message) {
	data, err := proto.Marshal(resp)
	if err != nil {
		c.String(http.StatusInternalServerError, "marshal resp error")
		return
	}
	c.Data(http.StatusOK, c.ContentType(), data)
	return
}

func CreateBotImpl(uin int64, password string) {
	log.Infof("开始读取设备信息")
	devicePath := fmt.Sprintf("device-%d.json", uin)
	bot.InitDevice(devicePath)
	log.Infof("设备信息 %+v", string(client.SystemDeviceInfo.ToJson()))

	log.Infof("创建机器人 %+v", uin)
	bot.Cli = client.NewClient(uin, password)

	log.Infof("初始化日志")
	bot.InitLog(bot.Cli)

	log.Infof("加载日志插件 Log")
	plugin.AddPrivateMessagePlugin(plugins.LogPrivateMessage)
	plugin.AddGroupMessagePlugin(plugins.LogGroupMessage)

	log.Infof("加载测试插件 Hello")
	plugin.AddPrivateMessagePlugin(plugins.HelloPrivateMessage)

	log.Infof("加载上报插件 Report")
	plugin.AddPrivateMessagePlugin(plugins.ReportPrivateMessage)
	plugin.AddGroupMessagePlugin(plugins.ReportGroupMessage)
	plugin.AddTempMessagePlugin(plugins.ReportTempMessage)
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

	plugin.Serve(bot.Cli)
	log.Infof("插件加载完成")

	log.Infof("登录中...")
	ok, err := bot.Login(bot.Cli)
	if err != nil {
		util.FatalError(fmt.Errorf("failed to login, err: %+v", err))
	}
	if ok {
		log.Infof("登录成功")
	} else {
		log.Infof("登录失败")
	}

	for {
		time.Sleep(5 * time.Second)
		if bot.Cli.Online {
			break
		}
		log.Infof("机器人不在线，可能在等待输入验证码")
	}

	log.Infof("刷新好友列表")
	if err := bot.Cli.ReloadFriendList(); err != nil {
		util.FatalError(fmt.Errorf("failed to load friend list, err: %+v", err))
	}
	log.Infof("共加载 %v 个好友.", len(bot.Cli.FriendList))

	log.Infof("刷新群列表")
	if err := bot.Cli.ReloadGroupList(); err != nil {
		util.FatalError(fmt.Errorf("failed to load group list, err: %+v", err))
	}
	log.Infof("共加载 %v 个群.", len(bot.Cli.GroupList))

	bot.ConnectUniversal(bot.Cli)
	go bot.Ping()
	go bot.ListenApi(bot.Cli)

	bot.SetRelogin(bot.Cli, 10, 30)
}
