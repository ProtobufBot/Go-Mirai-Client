package bot

import (
	"io"
	"net/http"

	"os"
	"strconv"
	"strings"
	"time"

	"github.com/2mf8/Go-Lagrange-Client/pkg/cache"
	"github.com/2mf8/Go-Lagrange-Client/pkg/util"
	"github.com/2mf8/Go-Lagrange-Client/proto_gen/onebot"

	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/message"
	log "github.com/sirupsen/logrus"
)

func EmptyText() *message.TextElement {
	return message.NewText("")
}

// 消息列表，不自动把code变成msg
func ProtoMsgToMiraiMsg(cli *client.QQClient, msgList []*onebot.Message, notConvertText bool) []message.IMessageElement {
	containReply := false // 每条消息只能包含一个reply
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
				messageChain = append(messageChain, RawMsgToMiraiMsg(cli, text)...) // 转换xml码
			}
		case "at":
			messageChain = append(messageChain, ProtoAtToMiraiAt(protoMsg.Data))
		case "image":
			messageChain = append(messageChain, ProtoImageToMiraiImage(protoMsg.Data))
		case "img":
			messageChain = append(messageChain, ProtoImageToMiraiImage(protoMsg.Data))
		case "friend_image":
			messageChain = append(messageChain, ProtoPrivateImageToMiraiPrivateImage(protoMsg.Data))
		case "friend_img":
			messageChain = append(messageChain, ProtoPrivateImageToMiraiPrivateImage(protoMsg.Data))
		case "record":
			messageChain = append(messageChain, ProtoVoiceToMiraiVoice(protoMsg.Data))
		case "face":
			messageChain = append(messageChain, ProtoFaceToMiraiFace(protoMsg.Data))
		case "reply":
			if replyElement := ProtoReplyToMiraiReply(protoMsg.Data); replyElement != nil && !containReply {
				containReply = true
				messageChain = append([]message.IMessageElement{replyElement}, messageChain...)
			}
		case "sleep":
			ProtoSleep(protoMsg.Data)
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
	elem := &message.GroupImageElement{}
	url, ok := data["url"]
	if !ok {
		url, ok = data["src"] // TODO 为了兼容我的旧代码偷偷加的
		if !ok {
			url, ok = data["file"]
		}
	}
	if !ok {
		log.Warnf("imageUrl不存在")
		return EmptyText()
	}
	b, err := preprocessImageMessage(url)
	if err == nil {
		elem.Stream = b
	}
	return elem
}

func ProtoPrivateImageToMiraiPrivateImage(data map[string]string) message.IMessageElement {
	elem := &message.FriendImageElement{}
	url, ok := data["url"]
	if !ok {
		url, ok = data["src"] // TODO 为了兼容我的旧代码偷偷加的
		if !ok {
			url, ok = data["file"]
		}
	}
	if !ok {
		log.Warnf("imageUrl不存在")
		return EmptyText()
	}
	b, err := preprocessImageMessage(url)
	if err == nil {
		elem.Stream = b
	}
	return elem
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
		return message.NewAt(0)
	}
	userId, err := strconv.ParseInt(qq, 10, 64)
	if err != nil {
		log.Warnf("atQQ不是数字")
		return EmptyText()
	}
	return message.NewAt(uint32(userId))
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
	return &message.FaceElement{
		FaceID: uint16(id),
	}
}

func ProtoReplyToMiraiReply(data map[string]string) *message.ReplyElement {
	rawMessage, hasRawMessage := data["raw_message"] // 如果存在 raw_message，按照raw_message显示

	messageIdStr, ok := data["message_id"]
	if !ok {
		return nil
	}
	messageIdInt, err := strconv.Atoi(messageIdStr)
	if err != nil {
		return nil
	}
	messageId := int32(messageIdInt)
	eventInterface, ok := cache.GroupMessageLru.Get(messageId)
	if ok {
		groupMessage, ok := eventInterface.(*message.GroupMessage)
		if ok {
			return &message.ReplyElement{
				ReplySeq: groupMessage.Id,
				Sender:   uint64(groupMessage.Sender.Uin),
				Time:     int32(groupMessage.Time),
				Elements: func() []message.IMessageElement {
					if hasRawMessage {
						return []message.IMessageElement{message.NewText(rawMessage)}
					} else {
						return groupMessage.Elements
					}
				}(),
			}
		}
	}
	eventInterface, ok = cache.PrivateMessageLru.Get(messageId)
	if ok {
		privateMessage, ok := eventInterface.(*message.PrivateMessage)
		if ok {
			return &message.ReplyElement{
				ReplySeq: privateMessage.Id,
				Sender:   uint64(privateMessage.Sender.Uin),
				Time:     privateMessage.Time,
				Elements: func() []message.IMessageElement {
					if hasRawMessage {
						return []message.IMessageElement{message.NewText(rawMessage)}
					} else {
						return privateMessage.Elements
					}
				}(),
			}
		}
	}
	return nil
}

func ProtoSleep(data map[string]string) {
	t, ok := data["time"]
	if !ok {
		log.Warnf("failed to get sleep time1")
		return
	}
	ms, err := strconv.Atoi(t)
	if err != nil {
		log.Warnf("failed to get sleep time2, %+v", err)
		return
	}
	if ms > 24*3600*1000 {
		log.Warnf("最多 sleep 24小时")
		ms = 24 * 3600 * 1000
	}
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
func preprocessImageMessage(path string) ([]byte, error) {
	if strings.Contains(path, "http") {
		resp, err := http.Get(path)
		defer resp.Body.Close()
		if err != nil {
			return nil, err
		}
		imo, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return imo, nil
	} else {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		reader, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}
		return reader, nil
	}
}

/*func ProtoMusicToMiraiMusic(_ *client.QQClient, data map[string]string) (m message.IMessageElement) {
	if data["type"] == "qq" {
		info, err := util.QQMusicSongInfo(data["id"])
		if err != nil {
			log.Warnf("failed to get qq music song info, %+v", data["id"])
			return EmptyText()
		}
		if !info.Get("track_info").Exists() {
			log.Warnf("music track_info not found, %+v", info.String())
			return EmptyText()
		}
		name := info.Get("track_info.name").Str
		mid := info.Get("track_info.mid").Str
		albumMid := info.Get("track_info.album.mid").Str
		pinfo, _ := util.GetBytes("http://u.y.qq.com/cgi-bin/musicu.fcg?g_tk=2034008533&uin=0&format=json&data={\"comm\":{\"ct\":23,\"cv\":0},\"url_mid\":{\"module\":\"vkey.GetVkeyServer\",\"method\":\"CgiGetVkey\",\"param\":{\"guid\":\"4311206557\",\"songmid\":[\"" + mid + "\"],\"songtype\":[0],\"uin\":\"0\",\"loginflag\":1,\"platform\":\"23\"}}}&_=1599039471576")
		jumpURL := "https://i.y.qq.com/v8/playsong.html?platform=11&appshare=android_qq&appversion=10030010&hosteuin=oKnlNenz7i-s7c**&songmid=" + mid + "&type=0&appsongtype=1&_wv=1&source=qq&ADTAG=qfshare"
		purl := gjson.ParseBytes(pinfo).Get("url_mid.data.midurlinfo.0.purl").Str
		preview := "http://y.gtimg.cn/music/photo_new/T002R180x180M000" + albumMid + ".jpg"
		content := info.Get("track_info.singer.0.name").Str
		if data["content"] != "" {
			content = data["content"]
		}
		return &message.MusicShareElement{
			MusicType:  message.QQMusic,
			Title:      name,
			Summary:    content,
			Url:        jumpURL,
			PictureUrl: preview,
			MusicUrl:   purl,
		}
	}
	if data["type"] == "163" {
		info, err := util.NeteaseMusicSongInfo(data["id"])
		if err != nil {
			log.Warnf("failed to get qq music song info, %+v", data["id"])
			return EmptyText()
		}
		if !info.Exists() {
			log.Warnf("netease song not fount")
			return EmptyText()
		}
		name := info.Get("name").Str
		jumpURL := "https://y.music.163.com/m/song/" + data["id"]
		musicURL := "http://music.163.com/song/media/outer/url?id=" + data["id"]
		picURL := info.Get("album.picUrl").Str
		artistName := ""
		if info.Get("artists.0").Exists() {
			artistName = info.Get("artists.0.name").Str
		}
		return &message.MusicShareElement{
			MusicType:  message.CloudMusic,
			Title:      name,
			Summary:    artistName,
			Url:        jumpURL,
			PictureUrl: picURL,
			MusicUrl:   musicURL,
		}
	}
	if data["type"] == "custom" {
		if data["subtype"] != "" {
			var subType int
			switch data["subtype"] {
			default:
				subType = message.QQMusic
			case "163":
				subType = message.CloudMusic
			case "migu":
				subType = message.MiguMusic
			case "kugou":
				subType = message.KugouMusic
			case "kuwo":
				subType = message.KuwoMusic
			}
			return &message.MusicShareElement{
				MusicType:  subType,
				Title:      data["title"],
				Summary:    data["content"],
				Url:        data["url"],
				PictureUrl: data["image"],
				MusicUrl:   data["audio"],
			}
		}
		xml := fmt.Sprintf(`<?xml version='1.0' encoding='UTF-8' standalone='yes' ?><msg serviceID="2" templateID="1" action="web" brief="[分享] %s" sourceMsgId="0" url="%s" flag="0" adverSign="0" multiMsgFlag="0"><item layout="2"><audio cover="%s" src="%s"/><title>%s</title><summary>%s</summary></item><source name="音乐" icon="https://i.gtimg.cn/open/app_icon/01/07/98/56/1101079856_100_m.png" url="http://web.p.qq.com/qqmpmobile/aio/app.html?id=1101079856" action="app" a_actionData="com.tencent.qqmusic" i_actionData="tencent1101079856://" appid="1101079856" /></msg>`,
			utils.XmlEscape(data["title"]), data["url"], data["image"], data["audio"], utils.XmlEscape(data["title"]), utils.XmlEscape(data["content"]))
		return &message.ServiceElement{
			Id:      60,
			Content: xml,
			SubType: "music",
		}
	}
	return EmptyText()
}*/
