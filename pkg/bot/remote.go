package bot

import (
	"bytes"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/ProtobufBot/Go-Mirai-Client/pkg/config"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/safe_ws"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/util"
	"github.com/ProtobufBot/Go-Mirai-Client/proto_gen/onebot"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

//go:generate go run github.com/a8m/syncmap -o "gen_remote_map.go" -pkg bot -name RemoteMap "map[int64]map[string]*WsServer"
var (
	// RemoteServers key是botId，value是map（key是serverName，value是server）
	RemoteServers RemoteMap
	jsonMarshaler = jsonpb.Marshaler{
		OrigName:     true,
		EmitDefaults: true,
	}
	jsonUnmarshaler = jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}
)

type WsServer struct {
	*safe_ws.SafeWebSocket        // 线程安全的ws
	*config.Plugin                // 服务器组配置
	wsUrl                  string // 随机抽中的url
	regexp                 *regexp.Regexp
}

func ConnectUniversal(cli *client.QQClient) {
	botServers := map[string]*WsServer{}
	RemoteServers.Store(cli.Uin, botServers)

	plugins := make([]*config.Plugin, 0)
	config.Plugins.Range(func(key string, value *config.Plugin) bool {
		plugins = append(plugins, value)
		return true
	})

	for _, group := range plugins {
		if group.Disabled || group.Urls == nil || len(group.Urls) < 1 {
			continue
		}
		serverGroup := *group
		util.SafeGo(func() {
			rand.Shuffle(len(serverGroup.Urls), func(i, j int) { serverGroup.Urls[i], serverGroup.Urls[j] = serverGroup.Urls[j], serverGroup.Urls[i] })
			urlIndex := 0 // 使用第几个url
			for IsClientExist(cli.Uin) {
				urlIndex = (urlIndex + 1) % len(serverGroup.Urls)
				serverUrl := serverGroup.Urls[urlIndex]
				log.Infof("开始连接Websocket服务器 [%s](%s)", serverGroup.Name, serverUrl)
				header := http.Header{}
				for k, v := range serverGroup.ExtraHeader {
					if v != nil {
						header[k] = v
					}
				}
				header["X-Self-ID"] = []string{strconv.FormatInt(cli.Uin, 10)}
				header["X-Client-Role"] = []string{"Universal"}
				conn, _, err := websocket.DefaultDialer.Dial(serverUrl, header)
				if err != nil {
					log.Warnf("连接Websocket服务器 [%s](%s) 错误，5秒后重连: %v", serverGroup.Name, serverUrl, err)
					time.Sleep(5 * time.Second)
					continue
				}
				log.Infof("连接Websocket服务器成功 [%s](%s)", serverGroup.Name, serverUrl)
				closeChan := make(chan int, 1)
				safeWs := safe_ws.NewSafeWebSocket(conn, OnWsRecvMessage(cli, &serverGroup), func() {
					defer func() {
						_ = recover() // 可能多次触发
					}()
					closeChan <- 1
				})
				botServers[serverGroup.Name] = &WsServer{
					SafeWebSocket: safeWs,
					Plugin:        &serverGroup,
					wsUrl:         serverUrl,
					regexp:        nil,
				}
				if serverGroup.RegexFilter != "" {
					if regex, err := regexp.Compile(serverGroup.RegexFilter); err != nil {
						log.Errorf("failed to compile [%s], regex_filter: %s", serverGroup.Name, serverGroup.RegexFilter)
					} else {
						botServers[serverGroup.Name].regexp = regex
					}
				}
				util.SafeGo(func() {
					for IsClientExist(cli.Uin) {
						if err := safeWs.Send(websocket.PingMessage, []byte("ping")); err != nil {
							break
						}
						time.Sleep(5 * time.Second)
					}
				})
				<-closeChan
				close(closeChan)
				delete(botServers, serverGroup.Name)
				log.Warnf("Websocket 服务器 [%s](%s) 已断开，5秒后重连", serverGroup.Name, serverUrl)
				time.Sleep(5 * time.Second)
			}
			log.Errorf("client does not exist, close websocket, %+v", cli.Uin)
		})
	}
}

func OnWsRecvMessage(cli *client.QQClient, plugin *config.Plugin) func(ws *safe_ws.SafeWebSocket, messageType int, data []byte) {
	apiFilter := map[onebot.Frame_FrameType]bool{}
	for _, apiType := range plugin.ApiFilter {
		apiFilter[onebot.Frame_FrameType(apiType)] = true
	}
	isApiAllow := func(frameType onebot.Frame_FrameType) bool {
		if len(apiFilter) == 0 {
			return true
		}
		return apiFilter[frameType]
	}

	return func(ws *safe_ws.SafeWebSocket, messageType int, data []byte) {
		if !IsClientExist(cli.Uin) {
			ws.Close()
			return
		}
		if messageType == websocket.PingMessage || messageType == websocket.PongMessage {
			return
		}
		if !cli.Online.Load() {
			log.Warnf("bot is not online, ignore API, %+v", cli.Uin)
			return
		}
		var apiReq onebot.Frame
		switch messageType {
		case websocket.BinaryMessage:
			err := proto.Unmarshal(data, &apiReq)
			if err != nil {
				log.Errorf("收到API binary，解析错误 %v", err)
				return
			}
		case websocket.TextMessage:
			err := jsonUnmarshaler.Unmarshal(bytes.NewReader(data), &apiReq)
			if err != nil {
				log.Errorf("收到API text，解析错误 %v", err)
				return
			}
		}

		log.Debugf("收到 apiReq 信息, %+v", util.MustMarshal(apiReq))

		apiResp := handleApiFrame(cli, &apiReq, isApiAllow)
		var (
			respBytes []byte
			err       error
		)
		switch messageType {
		case websocket.BinaryMessage:
			respBytes, err = proto.Marshal(apiResp)
			if err != nil {
				log.Errorf("failed to marshal api resp, %+v", err)
			}
		case websocket.TextMessage:
			respStr, err := jsonMarshaler.MarshalToString(apiResp)
			if err != nil {
				log.Errorf("failed to marshal api resp, %+v", err)
			}
			respBytes = []byte(respStr)
		}
		log.Debugf("发送 apiResp 信息, %+v", util.MustMarshal(apiResp))
		_ = ws.Send(messageType, respBytes)
	}
}

func handleApiFrame(cli *client.QQClient, req *onebot.Frame, isApiAllow func(onebot.Frame_FrameType) bool) (resp *onebot.Frame) {
	resp = &onebot.Frame{
		BotId: cli.Uin,
		Echo:  req.Echo,
		Ok:    true,
	}
	switch data := req.Data.(type) {
	case *onebot.Frame_SendPrivateMsgReq:
		resp.FrameType = onebot.Frame_TSendPrivateMsgResp
		if resp.Ok = isApiAllow(onebot.Frame_TSendPrivateMsgReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_SendPrivateMsgResp{
			SendPrivateMsgResp: HandleSendPrivateMsg(cli, data.SendPrivateMsgReq),
		}
	case *onebot.Frame_SendGroupMsgReq:
		resp.FrameType = onebot.Frame_TSendGroupMsgResp
		if resp.Ok = isApiAllow(onebot.Frame_TSendGroupMsgReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_SendGroupMsgResp{
			SendGroupMsgResp: HandleSendGroupMsg(cli, data.SendGroupMsgReq),
		}
	case *onebot.Frame_SendMsgReq:
		resp.FrameType = onebot.Frame_TSendMsgResp
		if resp.Ok = isApiAllow(onebot.Frame_TSendMsgReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_SendMsgResp{
			SendMsgResp: HandleSendMsg(cli, data.SendMsgReq),
		}
	case *onebot.Frame_DeleteMsgReq:
		resp.FrameType = onebot.Frame_TDeleteMsgResp
		if resp.Ok = isApiAllow(onebot.Frame_TDeleteMsgReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_DeleteMsgResp{
			DeleteMsgResp: HandleDeleteMsg(cli, data.DeleteMsgReq),
		}
	case *onebot.Frame_GetMsgReq:
		resp.FrameType = onebot.Frame_TGetMsgResp
		if resp.Ok = isApiAllow(onebot.Frame_TGetMsgReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_GetMsgResp{
			GetMsgResp: HandleGetMsg(cli, data.GetMsgReq),
		}
	case *onebot.Frame_SetGroupKickReq:
		resp.FrameType = onebot.Frame_TSetGroupKickResp
		if resp.Ok = isApiAllow(onebot.Frame_TSetGroupKickReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_SetGroupKickResp{
			SetGroupKickResp: HandleSetGroupKick(cli, data.SetGroupKickReq),
		}
	case *onebot.Frame_SetGroupBanReq:
		resp.FrameType = onebot.Frame_TSetGroupBanResp
		if resp.Ok = isApiAllow(onebot.Frame_TSetGroupBanReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_SetGroupBanResp{
			SetGroupBanResp: HandleSetGroupBan(cli, data.SetGroupBanReq),
		}
	case *onebot.Frame_SetGroupWholeBanReq:
		resp.FrameType = onebot.Frame_TSetGroupWholeBanResp
		if resp.Ok = isApiAllow(onebot.Frame_TSetGroupWholeBanReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_SetGroupWholeBanResp{
			SetGroupWholeBanResp: HandleSetGroupWholeBan(cli, data.SetGroupWholeBanReq),
		}
	case *onebot.Frame_SetGroupCardReq:
		resp.FrameType = onebot.Frame_TSetGroupCardResp
		if resp.Ok = isApiAllow(onebot.Frame_TSetGroupCardReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_SetGroupCardResp{
			SetGroupCardResp: HandleSetGroupCard(cli, data.SetGroupCardReq),
		}
	case *onebot.Frame_SetGroupNameReq:
		resp.FrameType = onebot.Frame_TSetGroupNameResp
		if resp.Ok = isApiAllow(onebot.Frame_TSetGroupNameReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_SetGroupNameResp{
			SetGroupNameResp: HandleSetGroupName(cli, data.SetGroupNameReq),
		}
	case *onebot.Frame_SetGroupLeaveReq:
		resp.FrameType = onebot.Frame_TSetGroupLeaveResp
		if resp.Ok = isApiAllow(onebot.Frame_TSetGroupLeaveReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_SetGroupLeaveResp{
			SetGroupLeaveResp: HandleSetGroupLeave(cli, data.SetGroupLeaveReq),
		}
	case *onebot.Frame_SetGroupSpecialTitleReq:
		resp.FrameType = onebot.Frame_TSetGroupSpecialTitleResp
		if resp.Ok = isApiAllow(onebot.Frame_TSetGroupSpecialTitleReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_SetGroupSpecialTitleResp{
			SetGroupSpecialTitleResp: HandleSetGroupSpecialTitle(cli, data.SetGroupSpecialTitleReq),
		}
	case *onebot.Frame_SetFriendAddRequestReq:
		resp.FrameType = onebot.Frame_TSetFriendAddRequestResp
		if resp.Ok = isApiAllow(onebot.Frame_TSetFriendAddRequestReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_SetFriendAddRequestResp{
			SetFriendAddRequestResp: HandleSetFriendAddRequest(cli, data.SetFriendAddRequestReq),
		}
	case *onebot.Frame_SetGroupAddRequestReq:
		resp.FrameType = onebot.Frame_TSetGroupAddRequestResp
		if resp.Ok = isApiAllow(onebot.Frame_TSetGroupAddRequestReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_SetGroupAddRequestResp{
			SetGroupAddRequestResp: HandleSetGroupAddRequest(cli, data.SetGroupAddRequestReq),
		}
	case *onebot.Frame_GetLoginInfoReq:
		resp.FrameType = onebot.Frame_TGetLoginInfoResp
		if resp.Ok = isApiAllow(onebot.Frame_TGetLoginInfoReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_GetLoginInfoResp{
			GetLoginInfoResp: HandleGetLoginInfo(cli, data.GetLoginInfoReq),
		}
	case *onebot.Frame_GetFriendListReq:
		resp.FrameType = onebot.Frame_TGetFriendListResp
		if resp.Ok = isApiAllow(onebot.Frame_TGetFriendListReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_GetFriendListResp{
			GetFriendListResp: HandleGetFriendList(cli, data.GetFriendListReq),
		}
	case *onebot.Frame_GetGroupInfoReq:
		resp.FrameType = onebot.Frame_TGetGroupInfoResp
		if resp.Ok = isApiAllow(onebot.Frame_TGetGroupInfoReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_GetGroupInfoResp{
			GetGroupInfoResp: HandleGetGroupInfo(cli, data.GetGroupInfoReq),
		}
	case *onebot.Frame_GetGroupListReq:
		resp.FrameType = onebot.Frame_TGetGroupListResp
		if resp.Ok = isApiAllow(onebot.Frame_TGetGroupListReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_GetGroupListResp{
			GetGroupListResp: HandleGetGroupList(cli, data.GetGroupListReq),
		}
	case *onebot.Frame_GetGroupMemberInfoReq:
		resp.FrameType = onebot.Frame_TGetGroupMemberInfoResp
		if resp.Ok = isApiAllow(onebot.Frame_TGetGroupMemberInfoReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_GetGroupMemberInfoResp{
			GetGroupMemberInfoResp: HandleGetGroupMemberInfo(cli, data.GetGroupMemberInfoReq),
		}
	case *onebot.Frame_GetGroupMemberListReq:
		resp.FrameType = onebot.Frame_TGetGroupMemberListResp
		if resp.Ok = isApiAllow(onebot.Frame_TGetGroupMemberListReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_GetGroupMemberListResp{
			GetGroupMemberListResp: HandleGetGroupMemberList(cli, data.GetGroupMemberListReq),
		}
	case *onebot.Frame_GetStrangerInfoReq:
		resp.FrameType = onebot.Frame_TGetStrangerInfoResp
		if resp.Ok = isApiAllow(onebot.Frame_TGetStrangerInfoReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_GetStrangerInfoResp{
			GetStrangerInfoResp: HandleGetStrangerInfo(cli, data.GetStrangerInfoReq),
		}
	case *onebot.Frame_GetCookiesReq:
		resp.FrameType = onebot.Frame_TGetCookiesResp
		if resp.Ok = isApiAllow(onebot.Frame_TGetCookiesReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_GetCookiesResp{
			GetCookiesResp: HandleGetCookies(cli, data.GetCookiesReq),
		}
	case *onebot.Frame_GetCsrfTokenReq:
		resp.FrameType = onebot.Frame_TGetCsrfTokenResp
		if resp.Ok = isApiAllow(onebot.Frame_TGetCsrfTokenReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_GetCsrfTokenResp{
			GetCsrfTokenResp: HandleGetCSRFToken(cli, data.GetCsrfTokenReq),
		}
	case *onebot.Frame_SetGroupSignInReq:
		resp.FrameType = onebot.Frame_TSetGroupSignInResp
		if resp.Ok = isApiAllow(onebot.Frame_TSetGroupSignInReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_SetGroupSignInResp{
			SetGroupSignInResp: HandleSetGroupSignIn(cli, data.SetGroupSignInReq),
		}
	case *onebot.Frame_SendMusicReq:
		resp.FrameType = onebot.Frame_TSendMusicResp
		if resp.Ok = isApiAllow(onebot.Frame_TSendMusicReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_SendMusicResp{
			SendMusicResp: HandleSendMusic(cli, data.SendMusicReq),
		}
	case *onebot.Frame_SendGroupPokeReq:
		resp.FrameType = onebot.Frame_TSendGroupPokeResp
		if resp.Ok = isApiAllow(onebot.Frame_TSendGroupPokeReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_SendGroupPokeResp{
			SendGroupPokeResp: HandleSendGroupPoke(cli, data.SendGroupPokeReq),
		}
	case *onebot.Frame_SendFriendPokeReq:
		resp.FrameType = onebot.Frame_TSendFriendPokeResp
		if resp.Ok = isApiAllow(onebot.Frame_TSendFriendPokeReq); !resp.Ok {
			return
		}
		resp.Data = &onebot.Frame_SendFriendPokeResp{
			SendFriendPokeResp: HandleSendFriendPoke(cli, data.SendFriendPokeReq),
		}
	default:
		return resp
	}
	return resp
}

func HandleEventFrame(cli *client.QQClient, eventFrame *onebot.Frame) {
	eventFrame.Ok = true
	eventFrame.BotId = cli.Uin
	eventBytes, err := proto.Marshal(eventFrame)
	if err != nil {
		log.Errorf("event 序列化错误 %v", err)
		return
	}

	wsServers, ok := RemoteServers.Load(cli.Uin)
	if !ok {
		log.Warnf("failed to load remote servers, %+v", cli.Uin)
		return
	}

	for _, ws := range wsServers {
		if ws.EventFilter != nil && len(ws.EventFilter) > 0 { // 有event filter
			if !int32SliceContains(ws.EventFilter, int32(eventFrame.FrameType)) {
				log.Debugf("EventFilter 跳过 [%s](%s)", ws.Name, ws.wsUrl)
				continue
			}
		}

		err := proto.Unmarshal(eventBytes, eventFrame) // 每个serverGroup, eventFrame 恢复原消息，防止因正则匹配互相影响
		if err != nil {
			log.Errorf("failed to unmarshal raw event frame, %+v", err)
			return
		}

		report := true // 是否上报event

		if ws.regexp != nil { // 有prefix filter
			if e, ok := eventFrame.Data.(*onebot.Frame_PrivateMessageEvent); ok {
				if report = ws.regexp.MatchString(e.PrivateMessageEvent.RawMessage); report && ws.RegexReplace != "" {
					e.PrivateMessageEvent.RawMessage = ws.regexp.ReplaceAllString(e.PrivateMessageEvent.RawMessage, ws.RegexReplace)
				}
			}
			if e, ok := eventFrame.Data.(*onebot.Frame_GroupMessageEvent); ok {
				if report = ws.regexp.MatchString(e.GroupMessageEvent.RawMessage); report && ws.RegexReplace != "" {
					e.GroupMessageEvent.RawMessage = ws.regexp.ReplaceAllString(e.GroupMessageEvent.RawMessage, ws.RegexReplace)
				}
			}
		}

		if report {
			if ws.Json {
				// 使用json上报
				sendingString, err := jsonMarshaler.MarshalToString(eventFrame)
				if err != nil {
					log.Errorf("event 序列化错误 %v", err)
					continue
				}
				_ = ws.Send(websocket.TextMessage, []byte(sendingString))
			} else {
				// 使用protobuf上报
				sendingBytes, err := proto.Marshal(eventFrame) // 使用正则修改后的eventFrame
				if err != nil {
					log.Errorf("event 序列化错误 %v", err)
					continue
				}
				log.Debugf("上报 event 给 [%s](%s)", ws.Name, ws.wsUrl)
				_ = ws.Send(websocket.BinaryMessage, sendingBytes)
			}
		}
	}
}

func int32SliceContains(numbers []int32, num int32) bool {
	for _, number := range numbers {
		if number == num {
			return true
		}
	}
	return false
}
