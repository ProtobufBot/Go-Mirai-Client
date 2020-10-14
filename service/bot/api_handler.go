package bot

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/ProtobufBot/Go-Mirai-Client/proto_gen/onebot"
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
	return &onebot.SendPrivateMsgResp{
		MessageId: ret.Id,
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
	return &onebot.SendGroupMsgResp{
		MessageId: ret.Id,
	}
}
