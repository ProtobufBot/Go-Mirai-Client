package bot

import (
	"strconv"
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/util"
	"github.com/ProtobufBot/Go-Mirai-Client/proto_gen/onebot"
	log "github.com/sirupsen/logrus"
)

func EmptyText() *message.TextElement {
	return message.NewText("")
}

// 消息列表，不自动把code变成msg
func ProtoMsgToMiraiMsg(msgList []*onebot.Message, notConvertText bool) []message.IMessageElement {
	messageChain := make([]message.IMessageElement, 0)
	for _, protoMsg := range msgList {
		switch protoMsg.Type {
		case "text":
			if notConvertText {
				messageChain = append(messageChain, ProtoTextToMiraiText(protoMsg.Data))
			} else {
				text, ok := protoMsg.Data["text"]
				if !ok {
					log.Warnf("text不存在")
					continue
				}
				messageChain = append(messageChain, RawMsgToMiraiMsg(text)...) // 转换xml码
			}
		case "at":
			messageChain = append(messageChain, ProtoAtToMiraiAt(protoMsg.Data))
		case "image":
			messageChain = append(messageChain, ProtoImageToMiraiImage(protoMsg.Data))
		case "img":
			messageChain = append(messageChain, ProtoImageToMiraiImage(protoMsg.Data))
		case "record":
			messageChain = append(messageChain, ProtoVoiceToMiraiVoice(protoMsg.Data))
		case "face":
			messageChain = append(messageChain, ProtoFaceToMiraiFace(protoMsg.Data))
		case "share":
			messageChain = append(messageChain, ProtoShareToMiraiShare(protoMsg.Data))
		case "light_app":
			messageChain = append(messageChain, ProtoLightAppToMiraiLightApp(protoMsg.Data))
		case "service":
			messageChain = append(messageChain, ProtoServiceToMiraiService(protoMsg.Data))
		default:
			log.Errorf("不支持的消息类型 %+v", protoMsg)
		}
	}
	return messageChain
}

func ProtoTextToMiraiText(data map[string]string) message.IMessageElement {
	text, ok := data["text"]
	if !ok {
		log.Warnf("text不存在")
		return EmptyText()
	}
	return message.NewText(text)
}

func ProtoImageToMiraiImage(data map[string]string) message.IMessageElement {
	url, ok := data["url"]
	if !ok || !strings.Contains(url, "http") {
		url, ok = data["src"] // TODO 为了兼容我的旧代码偷偷加的
		if !ok || !strings.Contains(url, "http") {
			url, ok = data["file"]
		}
	}
	if !ok || !strings.Contains(url, "http") {
		log.Warnf("imageUrl不存在")
		return EmptyText()
	}
	log.Infof("下载图片: %+v", url)
	b, err := util.GetBytes(url)
	if err != nil {
		log.Errorf("下载图片失败")
		return EmptyText()
	}
	return message.NewImage(b)
}

func ProtoVoiceToMiraiVoice(data map[string]string) message.IMessageElement {
	url, ok := data["url"]
	if !ok {
		url, ok = data["file"]
	}
	if !ok {
		log.Warnf("recordUrl不存在")
		return EmptyText()
	}
	b, err := util.GetBytes(url)
	if err != nil {
		log.Errorf("下载语音失败")
		return EmptyText()
	}
	if !util.IsAMRorSILK(b) {
		log.Errorf("不是amr或silk格式")
		return EmptyText()
	}
	return &message.VoiceElement{Data: b}
}

func ProtoAtToMiraiAt(data map[string]string) message.IMessageElement {
	qq, ok := data["qq"]
	if !ok {
		log.Warnf("atQQ不存在")
		return EmptyText()
	}
	if qq == "all" {
		return message.AtAll()
	}
	userId, err := strconv.ParseInt(qq, 10, 64)
	if err != nil {
		log.Warnf("atQQ不是数字")
		return EmptyText()
	}
	return message.NewAt(userId)
}

func ProtoFaceToMiraiFace(data map[string]string) message.IMessageElement {
	idStr, ok := data["id"]
	if !ok {
		log.Warnf("faceId不存在")
		return EmptyText()
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warnf("faceId不是数字")
		return EmptyText()
	}
	return message.NewFace(int32(id))
}

func ProtoShareToMiraiShare(data map[string]string) message.IMessageElement {
	url, ok := data["url"]
	if !ok {
		url = "https://www.baidu.com/"
	}
	title, ok := data["title"]
	if !ok {
		title = "分享标题"
	}
	content, ok := data["content"]
	if !ok {
		url = "分享内容"
	}
	image, ok := data["image"]
	if !ok {
		image = ""
	}
	return message.NewUrlShare(url, title, content, image)
}

func ProtoLightAppToMiraiLightApp(data map[string]string) message.IMessageElement {
	content, ok := data["content"]
	if !ok || content == "" {
		return EmptyText()
	}
	return message.NewLightApp(content)
}

func ProtoServiceToMiraiService(data map[string]string) message.IMessageElement {
	subType, ok := data["sub_type"]
	if !ok || subType == "" {
		log.Warnf("service sub_type不存在")
		return EmptyText()
	}

	content, ok := data["content"]
	if !ok {
		log.Warnf("service content不存在")
		return EmptyText()
	}

	id, ok := data["id"]
	if !ok {
		id = ""
	}
	resId, err := strconv.ParseInt(id, 10, 64)
	if err != nil || resId == 0 {
		if subType == "xml" {
			resId = 60 // xml默认60
		} else {
			resId = 1 // json默认1
		}
	}

	return &message.ServiceElement{
		Id:      int32(resId),
		Content: content,
		SubType: subType,
	}
}
