package bot

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/ProtobufBot/Go-Mirai-Client/proto_gen/onebot"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var Conn *websocket.Conn

var WsUrl = "ws://localhost:8081/ws/cq/"

var connecting = false
var connectLock sync.Mutex

func ConnectUniversal(cli *client.QQClient) {
	connectLock.Lock()
	if connecting {
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
			go listenApi(cli, conn)
			Conn = conn
			connectLock.Lock()
			connecting = false
			connectLock.Unlock()
			break
		}
	}
}

func listenApi(cli *client.QQClient, conn *websocket.Conn) {
	defer conn.Close()

	for {
		_, buf, err := conn.ReadMessage()
		if err != nil {
			log.Warnf("监听反向WS API时出现错误: %v", err)
			break
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
	case *onebot.Frame_DeleteMsgReq:
		resp.FrameType = onebot.Frame_TDeleteMsgReq
		resp.Data = &onebot.Frame_DeleteMsgResp{
			DeleteMsgResp: HandleDeleteMsg(cli, data.DeleteMsgReq),
		}
	case *onebot.Frame_GetMsgReq:
		resp.FrameType = onebot.Frame_TGetMsgReq
		resp.Data = &onebot.Frame_GetMsgResp{
			GetMsgResp: HandleGetMsg(cli, data.GetMsgReq),
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
