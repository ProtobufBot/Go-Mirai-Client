package bot

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/ProtobufBot/Go-Mirai-Client/proto_gen/onebot"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var Conn *websocket.Conn

var WsUrl = "ws://localhost:8081/ws/cq/"

var connecting = false
var connectLock sync.Mutex

func ConnectUniversal(cli *client.QQClient) {
	connectLock.Lock()
	if connecting {
		connectLock.Unlock()
		return
	}
	connecting = true
	connectLock.Unlock()
	header := http.Header{
		"X-Client-Role": []string{"Universal"},
		"X-Self-ID":     []string{strconv.FormatInt(cli.Uin, 10)},
		"User-Agent":    []string{"CQHttp/4.15.0"},
	}

	for {
		conn, _, err := websocket.DefaultDialer.Dial(WsUrl, header)
		if err != nil {
			log.Warnf("连接Websocket服务器 %v 时出现错误: %v", WsUrl, err)
			time.Sleep(5 * time.Second)
			continue
		} else {
			log.Infof("已连接Websocket %v", WsUrl)
			go ping(cli, conn)
			go listenApi(cli, conn)
			Conn = conn
			time.Sleep(1 * time.Second)
			connectLock.Lock()
			connecting = false
			connectLock.Unlock()
			break
		}
	}
}

func ping(cli *client.QQClient, conn *websocket.Conn) {
	errCount := 0
	for errCount < 5 {
		if err := conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
			log.Warnf("websocket ping失败")
			errCount++
		} else {
			errCount = 0
		}
		time.Sleep(10 * time.Second)
	}
	log.Warnf("websocket 连续ping失败5次，断开连接")
	_ = conn.Close()
	ConnectUniversal(cli)
}

func listenApi(cli *client.QQClient, conn *websocket.Conn) {
	defer conn.Close()

	for {
		messageType, buf, err := conn.ReadMessage()
		if err != nil {
			log.Warnf("监听反向WS API时出现错误: %v", err)
			break
		}
		if messageType == websocket.PingMessage || messageType == websocket.PongMessage {
			continue
		}

		var req onebot.Frame
		err = proto.Unmarshal(buf, &req)
		if err != nil {
			log.Errorf("收到API buffer，解析错误 %v", err)
			continue
		}

		go func() {
			resp := handleApiFrame(cli, &req)
			respBytes, err := resp.Marshal()
			if err != nil {
				log.Errorf("序列化ApiResp错误 %v", err)
			}
			err = conn.WriteMessage(websocket.BinaryMessage, respBytes)
			if err != nil {
				log.Errorf("发送ApiResp错误 %v", err)
			}
		}()
	}
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

	if Conn == nil {
		ConnectUniversal(cli)
		return
	}

	_ = Conn.SetWriteDeadline(time.Now().Add(time.Second * 15))
	err = Conn.WriteMessage(websocket.BinaryMessage, eventBytes)
	if err != nil {
		log.Errorf("发送Event错误 %v", err)
		_ = Conn.Close()
		ConnectUniversal(cli)
		return
	}
}
