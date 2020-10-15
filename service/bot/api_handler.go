package bot

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/ProtobufBot/Go-Mirai-Client/proto_gen/onebot"
	"github.com/ProtobufBot/Go-Mirai-Client/service/cache"
	log "github.com/sirupsen/logrus"
)

func HandleSendPrivateMsg(cli *client.QQClient, req *onebot.SendPrivateMsgReq) *onebot.SendPrivateMsgResp {
	miraiMsg := ProtoMsgToMiraiMsg(req.Message)

	var newElem []message.IMessageElement
	for _, elem := range miraiMsg {
		if i, ok := elem.(*message.ImageElement); ok {
			gm, err := cli.UploadPrivateImage(req.UserId, i.Data)
			if err != nil {
				log.Warnf("警告: 私聊图片上传失败: %v", err)
				continue
			}
			newElem = append(newElem, gm)
			continue
		}
		newElem = append(newElem, elem)
	}
	ret := cli.SendPrivateMessage(req.UserId, &message.SendingMessage{Elements: newElem})
	messageId := cache.NextGlobalSeq()
	cache.PrivateMessageLru.Add(messageId, ret)
	return &onebot.SendPrivateMsgResp{
		MessageId: messageId,
	}
}

func HandleSendGroupMsg(cli *client.QQClient, req *onebot.SendGroupMsgReq) *onebot.SendGroupMsgResp {
	miraiMsg := ProtoMsgToMiraiMsg(req.Message)
	var newElem []message.IMessageElement
	for _, elem := range miraiMsg {
		if i, ok := elem.(*message.ImageElement); ok {
			gm, err := cli.UploadGroupImage(req.GroupId, i.Data)
			if err != nil {
				log.Warnf("警告: 群聊图片上传失败: %v", err)
				continue
			}
			newElem = append(newElem, gm)
			continue
		}
		if i, ok := elem.(*message.VoiceElement); ok {
			gm, err := cli.UploadGroupPtt(req.GroupId, i.Data)
			if err != nil {
				log.Warnf("警告: 群聊语音上传失败: %v", err)
				continue
			}
			newElem = append(newElem, gm)
			continue
		}
		newElem = append(newElem, elem)
	}
	ret := cli.SendGroupMessage(req.GroupId, &message.SendingMessage{Elements: newElem})
	messageId := cache.NextGlobalSeq()
	cache.GroupMessageLru.Add(messageId, ret)
	return &onebot.SendGroupMsgResp{
		MessageId: messageId,
	}
}

func HandleDeleteMsg(cli *client.QQClient, req *onebot.DeleteMsgReq) *onebot.DeleteMsgResp {
	eventInterface, ok := cache.GroupMessageLru.Get(req.MessageId)
	if !ok {
		return &onebot.DeleteMsgResp{}
	}
	event := eventInterface.(*message.GroupMessage)
	cli.RecallGroupMessage(event.GroupCode, event.Id, event.InternalId)
	return &onebot.DeleteMsgResp{}
}
