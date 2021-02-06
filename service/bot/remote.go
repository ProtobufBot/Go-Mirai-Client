package bot

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/ProtobufBot/Go-Mirai-Client/config"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/safe_ws"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/util"
	"github.com/ProtobufBot/Go-Mirai-Client/proto_gen/onebot"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	WsServers = make(map[string]*safe_ws.SafeWebSocket) // TODO 线程安全？改用sync.map
)

func ConnectUniversal(cli *client.QQClient) {
	header := http.Header{
		"X-Client-Role": []string{"Universal"},
		"X-Self-ID":     []string{strconv.FormatInt(cli.Uin, 10)},
		"User-Agent":    []string{"CQHttp/4.15.0"},
	}
	for _, group := range config.Conf.ServerGroups {
		if group.Disabled {
			continue
		}
		serverGroup := *group
		util.SafeGo(func() {
			for {
				serverUrl := serverGroup.Urls[rand.Intn(len(serverGroup.Urls))]
				log.Infof("开始连接Websocket服务器 [%s](%s)", serverGroup.Name, serverUrl)
				conn, _, err := websocket.DefaultDialer.Dial(serverUrl, header)
				if err != nil {
					log.Warnf("连接Websocket服务器 [%s](%s) 错误，5秒后重连: %v", serverGroup.Name, serverUrl, err)
					time.Sleep(5 * time.Second)
					continue
				}
				log.Infof("连接Websocket服务器成功 [%s](%s)", serverGroup.Name, serverUrl)
				closeChan := make(chan int, 1)
				safeWs := safe_ws.NewSafeWebSocket(conn, OnWsRecvMessage, func() {
					defer func() {
						_ = recover() // 可能多次触发
					}()
					closeChan <- 1
				})
				WsServers[serverGroup.Name] = safeWs
				util.SafeGo(func() {
					for {
						if err := safeWs.Send(websocket.PingMessage, []byte("ping")); err != nil {
							break
						}
						time.Sleep(5 * time.Second)
					}
				})
				<-closeChan
				close(closeChan)
				delete(WsServers, serverGroup.Name)
				log.Warnf("Websocket 服务器 [%s](%s) 已断开，5秒后重连", serverGroup.Name, serverUrl)
				time.Sleep(5 * time.Second)
			}
		})
	}
}

func OnWsRecvMessage(ws *safe_ws.SafeWebSocket, messageType int, data []byte) {
	if messageType == websocket.PingMessage || messageType == websocket.PongMessage {
		return
	}
	var apiReq onebot.Frame
	err := proto.Unmarshal(data, &apiReq)
	if err != nil {
		log.Errorf("收到API buffer，解析错误 %v", err)
		return
	}
	log.Debugf("收到 apiReq 信息, %+v", util.MustMarshal(apiReq))

	apiResp := handleApiFrame(Cli, &apiReq)
	respBytes, err := apiResp.Marshal()
	if err != nil {
		log.Errorf("failed to marshal api resp, %+v", err)
	}
	log.Debugf("发送 apiResp 信息, %+v", util.MustMarshal(apiResp))
	_ = ws.Send(websocket.BinaryMessage, respBytes)
}

func handleApiFrame(cli *client.QQClient, req *onebot.Frame) *onebot.Frame {
	var resp = &onebot.Frame{
		BotId: cli.Uin,
		Echo:  req.Echo,
		Ok:    true,
	}
	switch data := req.Data.(type) {
	case *onebot.Frame_SendPrivateMsgReq:
		resp.FrameType = onebot.Frame_TSendPrivateMsgResp
		resp.Data = &onebot.Frame_SendPrivateMsgResp{
			SendPrivateMsgResp: HandleSendPrivateMsg(cli, data.SendPrivateMsgReq),
		}
	case *onebot.Frame_SendGroupMsgReq:
		resp.FrameType = onebot.Frame_TSendGroupMsgResp
		resp.Data = &onebot.Frame_SendGroupMsgResp{
			SendGroupMsgResp: HandleSendGroupMsg(cli, data.SendGroupMsgReq),
		}
	case *onebot.Frame_SendMsgReq:
		resp.FrameType = onebot.Frame_TSendMsgResp
		resp.Data = &onebot.Frame_SendMsgResp{
			SendMsgResp: HandleSendMsg(cli, data.SendMsgReq),
		}
	case *onebot.Frame_DeleteMsgReq:
		resp.FrameType = onebot.Frame_TDeleteMsgResp
		resp.Data = &onebot.Frame_DeleteMsgResp{
			DeleteMsgResp: HandleDeleteMsg(cli, data.DeleteMsgReq),
		}
	case *onebot.Frame_GetMsgReq:
		resp.FrameType = onebot.Frame_TGetMsgResp
		resp.Data = &onebot.Frame_GetMsgResp{
			GetMsgResp: HandleGetMsg(cli, data.GetMsgReq),
		}
	case *onebot.Frame_SetGroupKickReq:
		resp.FrameType = onebot.Frame_TSetGroupKickResp
		resp.Data = &onebot.Frame_SetGroupKickResp{
			SetGroupKickResp: HandleSetGroupKick(cli, data.SetGroupKickReq),
		}
	case *onebot.Frame_SetGroupBanReq:
		resp.FrameType = onebot.Frame_TSetGroupBanResp
		resp.Data = &onebot.Frame_SetGroupBanResp{
			SetGroupBanResp: HandleSetGroupBan(cli, data.SetGroupBanReq),
		}
	case *onebot.Frame_SetGroupWholeBanReq:
		resp.FrameType = onebot.Frame_TSetGroupWholeBanResp
		resp.Data = &onebot.Frame_SetGroupWholeBanResp{
			SetGroupWholeBanResp: HandleSetGroupWholeBan(cli, data.SetGroupWholeBanReq),
		}
	case *onebot.Frame_SetGroupCardReq:
		resp.FrameType = onebot.Frame_TSetGroupCardResp
		resp.Data = &onebot.Frame_SetGroupCardResp{
			SetGroupCardResp: HandleSetGroupCard(cli, data.SetGroupCardReq),
		}
	case *onebot.Frame_SetGroupNameReq:
		resp.FrameType = onebot.Frame_TSetGroupNameResp
		resp.Data = &onebot.Frame_SetGroupNameResp{
			SetGroupNameResp: HandleSetGroupName(cli, data.SetGroupNameReq),
		}
	case *onebot.Frame_SetGroupLeaveReq:
		resp.FrameType = onebot.Frame_TSetGroupLeaveResp
		resp.Data = &onebot.Frame_SetGroupLeaveResp{
			SetGroupLeaveResp: HandleSetGroupLeave(cli, data.SetGroupLeaveReq),
		}
	case *onebot.Frame_SetGroupSpecialTitleReq:
		resp.FrameType = onebot.Frame_TSetGroupSpecialTitleResp
		resp.Data = &onebot.Frame_SetGroupSpecialTitleResp{
			SetGroupSpecialTitleResp: HandleSetGroupSpecialTitle(cli, data.SetGroupSpecialTitleReq),
		}
	case *onebot.Frame_SetFriendAddRequestReq:
		resp.FrameType = onebot.Frame_TSetFriendAddRequestResp
		resp.Data = &onebot.Frame_SetFriendAddRequestResp{
			SetFriendAddRequestResp: HandleSetFriendAddRequest(cli, data.SetFriendAddRequestReq),
		}
	case *onebot.Frame_SetGroupAddRequestReq:
		resp.FrameType = onebot.Frame_TSetGroupAddRequestResp
		resp.Data = &onebot.Frame_SetGroupAddRequestResp{
			SetGroupAddRequestResp: HandleSetGroupAddRequest(cli, data.SetGroupAddRequestReq),
		}
	case *onebot.Frame_GetLoginInfoReq:
		resp.FrameType = onebot.Frame_TGetLoginInfoResp
		resp.Data = &onebot.Frame_GetLoginInfoResp{
			GetLoginInfoResp: HandleGetLoginInfo(cli, data.GetLoginInfoReq),
		}
	case *onebot.Frame_GetFriendListReq:
		resp.FrameType = onebot.Frame_TGetFriendListResp
		resp.Data = &onebot.Frame_GetFriendListResp{
			GetFriendListResp: HandleGetFriendList(cli, data.GetFriendListReq),
		}
	case *onebot.Frame_GetGroupInfoReq:
		resp.FrameType = onebot.Frame_TGetGroupInfoResp
		resp.Data = &onebot.Frame_GetGroupInfoResp{
			GetGroupInfoResp: HandleGetGroupInfo(cli, data.GetGroupInfoReq),
		}
	case *onebot.Frame_GetGroupListReq:
		resp.FrameType = onebot.Frame_TGetGroupListResp
		resp.Data = &onebot.Frame_GetGroupListResp{
			GetGroupListResp: HandleGetGroupList(cli, data.GetGroupListReq),
		}
	case *onebot.Frame_GetGroupMemberInfoReq:
		resp.FrameType = onebot.Frame_TGetGroupMemberInfoResp
		resp.Data = &onebot.Frame_GetGroupMemberInfoResp{
			GetGroupMemberInfoResp: HandleGetGroupMemberInfo(cli, data.GetGroupMemberInfoReq),
		}
	case *onebot.Frame_GetGroupMemberListReq:
		resp.FrameType = onebot.Frame_TGetGroupMemberListResp
		resp.Data = &onebot.Frame_GetGroupMemberListResp{
			GetGroupMemberListResp: HandleGetGroupMemberList(cli, data.GetGroupMemberListReq),
		}
	case *onebot.Frame_GetStrangerInfoReq:
		resp.FrameType = onebot.Frame_TGetStrangerInfoResp
		resp.Data = &onebot.Frame_GetStrangerInfoResp{
			GetStrangerInfoResp: HandleGetStrangerInfo(cli, data.GetStrangerInfoReq),
		}
	default:
		return resp
	}
	return resp
}

func HandleEventFrame(cli *client.QQClient, eventFrame *onebot.Frame) {
	eventFrame.Ok = true
	eventFrame.BotId = cli.Uin
	eventBytes, err := eventFrame.Marshal()
	if err != nil {
		log.Errorf("event 序列化错误 %v", err)
		return
	}

	for name, ws := range WsServers {
		log.Debugf("上报 event 给 [%s]", name)
		_ = ws.Send(websocket.BinaryMessage, eventBytes)
	}
}
