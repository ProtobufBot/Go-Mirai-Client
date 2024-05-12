package bot

import (
	"strconv"

	"github.com/2mf8/Go-Lagrange-Client/proto_gen/onebot"

	"github.com/2mf8/LagrangeGo/client"
	"github.com/2mf8/LagrangeGo/message"
)

func MiraiMsgToProtoMsg(cli *client.QQClient, messageChain []message.IMessageElement) []*onebot.Message {
	msgList := make([]*onebot.Message, 0)
	for _, element := range messageChain {
		switch elem := element.(type) {
		case *message.TextElement:
			msgList = append(msgList, MiraiTextToProtoText(elem))
		case *message.AtElement:
			msgList = append(msgList, MiraiAtToProtoAt(elem))
		case *message.FriendImageElement:
			msgList = append(msgList, MiraiFriendImageToProtoImage(elem))
		case *message.GroupImageElement:
			msgList = append(msgList, MiraiGroupImageToProtoImage(elem))
		case *message.FaceElement:
			msgList = append(msgList, MiraiFaceToProtoFace(elem))
		case *message.VoiceElement:
			msgList = append(msgList, MiraiVoiceToProtoVoice(elem))
		case *message.ShortVideoElement:
			msgList = append(msgList, MiraiVideoToProtoVideo(cli, elem))
		case *message.ReplyElement:
			msgList = append(msgList, MiraiReplyToProtoReply(cli, elem))
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

func MiraiFriendImageToProtoImage(elem *message.FriendImageElement) *onebot.Message {
	msg := &onebot.Message{
		Type: "image",
		Data: map[string]string{
			"image_id": elem.ImageId,
			"file":     elem.Url,
			"url":      elem.Url,
		},
	}
	if elem.Flash {
		msg.Data["type"] = "flash"
	}
	return msg
}

func MiraiGroupImageToProtoImage(elem *message.GroupImageElement) *onebot.Message {
	msg := &onebot.Message{
		Type: "image",
		Data: map[string]string{
			"image_id": elem.ImageId,
			"file":     elem.Url,
			"url":      elem.Url,
		},
	}
	if elem.Flash {
		msg.Data["type"] = "flash"
	}
	if elem.EffectID != 0 {
		msg.Data["type"] = "show"
		msg.Data["effect_id"] = strconv.FormatInt(int64(elem.EffectID), 10)
	}
	return msg
}

func MiraiAtToProtoAt(elem *message.AtElement) *onebot.Message {
	return &onebot.Message{
		Type: "at",
		Data: map[string]string{
			"qq": func() string {
				if elem.Target == 0 {
					return "all"
				}
				return strconv.FormatInt(int64(elem.Target), 10)
			}(),
		},
	}
}

func MiraiFaceToProtoFace(elem *message.FaceElement) *onebot.Message {
	return &onebot.Message{
		Type: "face",
		Data: map[string]string{
			"id": strconv.Itoa(int(elem.FaceID)),
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

func MiraiVideoToProtoVideo(cli *client.QQClient, elem *message.ShortVideoElement) *onebot.Message {
	return &onebot.Message{
		Type: "video",
		Data: map[string]string{
			"name": elem.Name,
			"url":  elem.Url,
		},
	}
}

func MiraiReplyToProtoReply(cli *client.QQClient, elem *message.ReplyElement) *onebot.Message {
	return &onebot.Message{
		Type: "reply",
		Data: map[string]string{
			"reply_seq":   strconv.FormatInt(int64(elem.ReplySeq), 10),
			"sender":      strconv.FormatInt(int64(elem.Sender), 10),
			"time":        strconv.FormatInt(int64(elem.Time), 10),
			"raw_message": MiraiMsgToRawMsg(cli, elem.Elements),
		},
	}
}