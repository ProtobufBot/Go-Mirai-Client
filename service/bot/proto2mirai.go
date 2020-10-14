package bot

import (
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/util"
	"github.com/ProtobufBot/Go-Mirai-Client/proto_gen/onebot"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func EmptyText() *message.TextElement {
	return message.NewText("")
}

func ProtoMsgToMiraiMsg(msgList []*onebot.Message) []message.IMessageElement {
	messageChain := make([]message.IMessageElement, 0)
	for _, protoMsg := range msgList {
		switch protoMsg.Type {
		case "at":
			messageChain = append(messageChain, ProtoAtToMiraiAt(protoMsg.Data))
		case "text":
			messageChain = append(messageChain, ProtoTextToMiraiText(protoMsg.Data))
		case "image":
			messageChain = append(messageChain, ProtoImageToMiraiImage(protoMsg.Data))
		case "face":
			messageChain = append(messageChain, ProtoFaceToMiraiFace(protoMsg.Data))
		default:
			log.Errorf("不支持的消息类型 %+v", protoMsg)

		}
	}
	return messageChain
}

func ProtoTextToMiraiText(data map[string]string) message.IMessageElement {
	text, ok := data["text"]
	if !ok {
		return EmptyText()
	}
	return message.NewText(text)
}

func ProtoImageToMiraiImage(data map[string]string) message.IMessageElement {
	url, ok := data["file"]
	if !ok {
		url, ok = data["url"]
	}
	if !ok {
		return EmptyText()
	}
	b, err := util.GetBytes(url)
	if err != nil {
		log.Errorf("获取图片失败")
		return EmptyText()
	}
	return message.NewImage(b)
}

func ProtoAtToMiraiAt(data map[string]string) message.IMessageElement {
	qq, ok := data["qq"]
	if !ok {
		return EmptyText()
	}
	if qq == "all" {
		return message.AtAll()
	}
	userId, err := strconv.ParseInt(qq, 10, 64)
	if err != nil {
		return EmptyText()
	}
	return message.NewAt(userId)
}

func ProtoFaceToMiraiFace(data map[string]string) message.IMessageElement {
	idStr, ok := data["id"]
	if !ok {
		return EmptyText()
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return EmptyText()
	}
	return message.NewFace(int32(id))
}
