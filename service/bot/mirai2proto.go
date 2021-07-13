package bot

import (
	"strconv"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/ProtobufBot/Go-Mirai-Client/proto_gen/onebot"
)

func MiraiMsgToProtoMsg(cli *client.QQClient, messageChain []message.IMessageElement) []*onebot.Message {
	msgList := make([]*onebot.Message, 0)
	for _, element := range messageChain {
		switch elem := element.(type) {
		case *message.TextElement:
			msgList = append(msgList, MiraiTextToProtoText(elem))
		case *message.AtElement:
			msgList = append(msgList, MiraiAtToProtoAt(elem))
		case *message.ImageElement:
			msgList = append(msgList, MiraiImageToProtoImage(elem))
		case *message.FaceElement:
			msgList = append(msgList, MiraiFaceToProtoFace(elem))
		case *message.VoiceElement:
			msgList = append(msgList, MiraiVoiceToProtoVoice(elem))
		case *message.ServiceElement:
			msgList = append(msgList, MiraiServiceToProtoService(elem))
		case *message.LightAppElement:
			msgList = append(msgList, MiraiLightAppToProtoLightApp(elem))
		case *message.ShortVideoElement:
			msgList = append(msgList, MiraiVideoToProtoVideo(cli, elem))
		case *message.ReplyElement:
			msgList = append(msgList, MiraiReplyToProtoReply(cli,elem))
		}
	}
	return msgList
}

func MiraiTextToProtoText(elem *message.TextElement) *onebot.Message {
	return &onebot.Message{
		Type: "text",
		Data: map[string]string{
			"text": elem.Content,
		},
	}
}

func MiraiImageToProtoImage(elem *message.ImageElement) *onebot.Message {
	return &onebot.Message{
		Type: "image",
		Data: map[string]string{
			"file": elem.Url,
			"url":  elem.Url,
		},
	}
}

func MiraiAtToProtoAt(elem *message.AtElement) *onebot.Message {
	return &onebot.Message{
		Type: "at",
		Data: map[string]string{
			"qq": func() string {
				if elem.Target == 0 {
					return "all"
				}
				return strconv.FormatInt(elem.Target, 10)
			}(),
		},
	}
}

func MiraiFaceToProtoFace(elem *message.FaceElement) *onebot.Message {
	return &onebot.Message{
		Type: "face",
		Data: map[string]string{
			"id": strconv.Itoa(int(elem.Index)),
		},
	}
}

func MiraiVoiceToProtoVoice(elem *message.VoiceElement) *onebot.Message {
	return &onebot.Message{
		Type: "record",
		Data: map[string]string{
			"file": elem.Url,
			"url":  elem.Url,
		},
	}
}

func MiraiServiceToProtoService(elem *message.ServiceElement) *onebot.Message {
	return &onebot.Message{
		Type: "service",
		Data: map[string]string{
			"id":       strconv.Itoa(int(elem.Id)),
			"content":  elem.Content,
			"res_id":   elem.ResId,
			"sub_type": elem.SubType,
		},
	}
}

func MiraiLightAppToProtoLightApp(elem *message.LightAppElement) *onebot.Message {
	return &onebot.Message{
		Type: "light_app",
		Data: map[string]string{
			"content": elem.Content,
		},
	}
}

func MiraiVideoToProtoVideo(cli *client.QQClient, elem *message.ShortVideoElement) *onebot.Message {
	return &onebot.Message{
		Type: "video",
		Data: map[string]string{
			"name": elem.Name,
			"url":  cli.GetShortVideoUrl(elem.Uuid, elem.Md5),
		},
	}
}

func MiraiReplyToProtoReply(cli *client.QQClient,elem *message.ReplyElement) *onebot.Message {
	return &onebot.Message{
		Type: "reply",
		Data: map[string]string{
			"reply_seq":   strconv.FormatInt(int64(elem.ReplySeq), 10),
			"sender":      strconv.FormatInt(elem.Sender, 10),
			"time":        strconv.FormatInt(int64(elem.Time), 10),
			"raw_message": MiraiMsgToRawMsg(cli,elem.Elements),
		},
	}
}
