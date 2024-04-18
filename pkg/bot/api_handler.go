package bot

import (
	_ "image/gif" // 用于解决发不出图片的问题
	_ "image/jpeg"
	_ "image/png"
	"math"
	"time"
	_ "unsafe"

	"github.com/ProtobufBot/Go-Mirai-Client/pkg/cache"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/config"
	"github.com/ProtobufBot/Go-Mirai-Client/proto_gen/onebot"

	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/message"
	log "github.com/sirupsen/logrus"
)

const MAX_TEXT_LENGTH = 80

// 风控临时解决方案
func splitText(content string, limit int) []string {
	text := []rune(content)

	result := make([]string, 0)
	num := int(math.Ceil(float64(len(text)) / float64(limit)))
	for i := 0; i < num; i++ {
		start := i * limit
		end := func() int {
			if (i+1)*limit > len(text) {
				return len(text)
			} else {
				return (i + 1) * limit
			}
		}()
		result = append(result, string(text[start:end]))
	}
	return result
}

func HandleSendPrivateMsg(cli *client.QQClient, req *onebot.SendPrivateMsgReq) *onebot.SendPrivateMsgResp {
	miraiMsg := ProtoMsgToMiraiMsg(cli, req.Message, req.AutoEscape)
	sendingMessage := &message.SendingMessage{Elements: miraiMsg}
	log.Infof("Bot(%d) Private(%d) <- %s", cli.Uin, req.UserId, MiraiMsgToRawMsg(cli, miraiMsg))
	ret, _ := cli.SendPrivateMessage(uint32(req.UserId), sendingMessage.Elements)
	cache.PrivateMessageLru.Add(ret.Result, ret)
	return &onebot.SendPrivateMsgResp{
		MessageId: ret.Result,
		MessageReceipt: &onebot.MessageReceipt{
			SenderId: req.UserId,
			Time:     time.Now().Unix(),
			Seqs:     []int32{ret.Result},
		},
	}
}

func HandleSendGroupMsg(cli *client.QQClient, req *onebot.SendGroupMsgReq) *onebot.SendGroupMsgResp {
	miraiMsg := ProtoMsgToMiraiMsg(cli, req.Message, req.AutoEscape)
	sendingMessage := &message.SendingMessage{Elements: miraiMsg}
	log.Infof("Bot(%d) Group(%d) <- %s", cli.Uin, req.GroupId, MiraiMsgToRawMsg(cli, miraiMsg))
	if len(sendingMessage.Elements) == 0 {
		log.Warnf("发送消息内容为空")
		return nil
	}
	ret, _ := cli.SendGroupMessage(uint32(req.GroupId), sendingMessage.Elements)
	if ret == nil || ret.Result == -1 {
		config.Fragment = !config.Fragment
		log.Warnf("发送群消息失败，可能被风控，下次发送将改变分片策略，Fragment: %+v", config.Fragment)
		return nil
	}
	cache.GroupMessageLru.Add(ret.Result, ret)
	return &onebot.SendGroupMsgResp{
		MessageId: ret.Result,
		MessageReceipt: &onebot.MessageReceipt{
			Time:    time.Now().Unix(),
			Seqs:    []int32{ret.Result},
			GroupId: req.GroupId,
		},
	}
}

func HandleSendMsg(cli *client.QQClient, req *onebot.SendMsgReq) *onebot.SendMsgResp {
	miraiMsg := ProtoMsgToMiraiMsg(cli, req.Message, req.AutoEscape)
	sendingMessage := &message.SendingMessage{Elements: miraiMsg}

	if req.GroupId != 0 && req.UserId != 0 { // 临时
		ret, _ := cli.SendTempMessage(uint32(req.GroupId), uint32(req.UserId), sendingMessage.Elements)
		cache.PrivateMessageLru.Add(ret.Result, ret)
		return &onebot.SendMsgResp{
			MessageId: ret.Result,
			MessageReceipt: &onebot.MessageReceipt{
				SenderId: req.UserId,
				Time:     time.Now().Unix(),
				Seqs:     []int32{ret.Result},
				GroupId:  req.GroupId,
			},
		}
	}

	if req.GroupId != 0 { // 群
		ret, _ := cli.SendGroupMessage(uint32(req.GroupId), sendingMessage.Elements)
		if ret == nil || ret.Result == -1 {
			config.Fragment = !config.Fragment
			log.Warnf("发送群消息失败，可能被风控，下次发送将改变分片策略，Fragment: %+v", config.Fragment)
			return nil
		}
		cache.GroupMessageLru.Add(ret.Result, ret)
		return &onebot.SendMsgResp{
			MessageId: ret.Result,
			MessageReceipt: &onebot.MessageReceipt{
				Time:    time.Now().Unix(),
				Seqs:    []int32{ret.Result},
				GroupId: req.GroupId,
			},
		}
	}

	if req.UserId != 0 { // 私聊
		ret, _ := cli.SendPrivateMessage(uint32(req.UserId), sendingMessage.Elements)
		cache.PrivateMessageLru.Add(ret.Result, ret)
		return &onebot.SendMsgResp{
			MessageId: ret.Result,
			MessageReceipt: &onebot.MessageReceipt{
				SenderId: req.UserId,
				Time:     time.Now().Unix(),
				Seqs:     []int32{ret.Result},
			},
		}
	}
	log.Warnf("failed to send msg")
	return nil
}

func HandleGetMsg(cli *client.QQClient, req *onebot.GetMsgReq) *onebot.GetMsgResp {
	eventInterface, isGroup := cache.GroupMessageLru.Get(req.MessageId)
	if isGroup {
		event := eventInterface.(*message.GroupMessage)
		messageType := "group"
		if event.Sender.Uin == cli.Uin {
			messageType = "self"
		}
		return &onebot.GetMsgResp{
			Time:        int32(event.Time),
			MessageType: messageType,
			MessageId:   req.MessageId,
			RealId:      event.InternalId, // 不知道是什么？
			Message:     MiraiMsgToProtoMsg(cli, event.Elements),
			RawMessage:  MiraiMsgToRawMsg(cli, event.Elements),
			Sender: &onebot.GetMsgResp_Sender{
				UserId:   int64(event.Sender.Uin),
				Nickname: event.Sender.Nickname,
			},
		}

	}
	eventInterface, isPrivate := cache.PrivateMessageLru.Get(req.MessageId)
	if isPrivate {
		event := eventInterface.(*message.PrivateMessage)
		messageType := "private"
		if event.Sender.Uin == cli.Uin {
			messageType = "self"
		}
		return &onebot.GetMsgResp{
			Time:        event.Time,
			MessageType: messageType,
			MessageId:   req.MessageId,
			RealId:      event.InternalId, // 不知道是什么？
			Message:     MiraiMsgToProtoMsg(cli, event.Elements),
			RawMessage:  MiraiMsgToRawMsg(cli, event.Elements),
			Sender: &onebot.GetMsgResp_Sender{
				UserId:   int64(event.Sender.Uin),
				Nickname: event.Sender.Nickname,
			},
		}
	}
	return nil
}