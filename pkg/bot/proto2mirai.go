package bot

import (
	"bytes"
	"fmt"

	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/ProtobufBot/Go-Mirai-Client/pkg/cache"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/clz"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/util"
	"github.com/ProtobufBot/Go-Mirai-Client/proto_gen/onebot"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/Mrs4s/MiraiGo/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
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
		case "dice":
			messageChain = append(messageChain, ProtoDiceToMiraiDice(protoMsg.Data))
		case "finger_guessing":
			messageChain = append(messageChain, ProtoFingerGuessingToMiraiFingerGuessing(protoMsg.Data))
		case "poke":
			messageChain = append(messageChain, ProtoPokeToMiraiPoke(protoMsg.Data))
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
		case "reply":
			if replyElement := ProtoReplyToMiraiReply(protoMsg.Data); replyElement != nil && !containReply {
				containReply = true
				messageChain = append([]message.IMessageElement{replyElement}, messageChain...)
			}
		case "sleep":
			ProtoSleep(protoMsg.Data)
		case "tts":
			messageChain = append(messageChain, ProtoTtsToMiraiTts(cli, protoMsg.Data))
		case "video":
			messageChain = append(messageChain, ProtoVideoToMiraiVideo(cli, protoMsg.Data))
		case "music":
			messageChain = append(messageChain, ProtoMusicToMiraiMusic(cli, protoMsg.Data))
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
	elem := &clz.LocalImageElement{}
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
	elem.Url = url
	if strings.Contains(url, "http://") || strings.Contains(url, "https://") {
		b, err := util.GetBytes(url)
		if err != nil {
			log.Errorf("failed to download image, %+v", err)
			return EmptyText()
		}
		elem.Stream = bytes.NewReader(b)
	} else {
		imageBytes, err := ioutil.ReadFile(url)
		if err != nil {
			log.Errorf("failed to open local image, %+v", err)
			return EmptyText()
		}
		elem.Stream = bytes.NewReader(imageBytes)
	}

	elem.Tp = data["type"] // show或flash
	if elem.Tp == "show" {
		effectIdStr := data["effect_id"]
		effectId, err := strconv.Atoi(effectIdStr)
		if err != nil || effectId < 40000 || effectId > 40005 {
			effectId = 40000
		}
		elem.EffectId = int32(effectId)
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
		return message.AtAll()
	}
	userId, err := strconv.ParseInt(qq, 10, 64)
	if err != nil {
		log.Warnf("atQQ不是数字")
		return EmptyText()
	}
	return message.NewAt(userId)
}

func ProtoDiceToMiraiDice(data map[string]string) message.IMessageElement {
	value := int32(rand.Intn(6) + 1)
	valueStr, ok := data["value"]
	if ok {
		v, err := strconv.ParseInt(valueStr, 10, 64)
		if err == nil && v >= 1 && v <= 6 {
			value = int32(v)
		}
	}
	return message.NewDice(value)
}

func ProtoFingerGuessingToMiraiFingerGuessing(data map[string]string) message.IMessageElement {
	value := int32(rand.Intn(2))
	valueStr, ok := data["value"]
	if ok {
		v, err := strconv.ParseInt(valueStr, 10, 64)
		if err == nil && v >= 0 && v <= 2 {
			value = int32(v)
		}
	}
	return message.NewFingerGuessing(value)
}

func ProtoPokeToMiraiPoke(data map[string]string) message.IMessageElement {
	qq, ok := data["qq"]
	if !ok {
		log.Warnf("pokeQQ不存在")
		return EmptyText()
	}
	userId, err := strconv.ParseInt(qq, 10, 64)
	if err != nil {
		log.Warnf("pokeQQ不是数字")
		return EmptyText()
	}
	return &clz.PokeElement{Target: userId}
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
				Sender:   groupMessage.Sender.Uin,
				Time:     groupMessage.Time,
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
				Sender:   privateMessage.Sender.Uin,
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

func ProtoTtsToMiraiTts(cli *client.QQClient, data map[string]string) (m message.IMessageElement) {
	defer func() {
		if r := recover(); r != nil {
			log.Warnf("tts 转换失败")
			m = EmptyText()
		}
	}()
	text, ok := data["text"]
	if !ok {
		log.Warnf("text不存在")
		return EmptyText()
	}
	b, err := cli.GetTts(text)
	if err != nil {
		log.Warnf("failed to get tts, %+v", err)
		return EmptyText()
	}
	return &message.VoiceElement{Data: b}
}

func ProtoVideoToMiraiVideo(_ *client.QQClient, data map[string]string) (m message.IMessageElement) {
	elem := &clz.MyVideoElement{}
	coverUrl, ok := data["cover"]
	if !ok {
		log.Warnf("video cover不存在")
		return EmptyText()
	}
	url, ok := data["url"]
	if !ok {
		url, ok = data["file"]
		if !ok {
			log.Warnf("video url不存在")
			return EmptyText()
		}
	}
	if strings.Contains(coverUrl, "http://") || strings.Contains(coverUrl, "https://") {
		coverBytes, err := util.GetBytes(coverUrl)
		if err != nil {
			log.Errorf("failed to download cover, err: %+v", err)
			return EmptyText()
		}
		elem.UploadingCover = bytes.NewReader(coverBytes)
	} else {
		coverBytes, err := ioutil.ReadFile(coverUrl)
		if err != nil {
			log.Errorf("failed to open file, err: %+v", err)
			return EmptyText()
		}
		elem.UploadingCover = bytes.NewReader(coverBytes)
	}

	videoFilePath := path.Join("video", util.MustMd5(url)+".mp4")
	if strings.Contains(url, "http://") || strings.Contains(url, "https://") {
		if !util.PathExists("video") {
			err := os.MkdirAll("video", 0777)
			if err != nil {
				log.Errorf("failed to mkdir, err: %+v", err)
				return EmptyText()
			}
		}
		if data["cache"] == "0" && util.PathExists(videoFilePath) {
			if err := os.Remove(videoFilePath); err != nil {
				log.Errorf("删除缓存文件 %v 时出现错误: %v", videoFilePath, err)
				return EmptyText()
			}
		}
		if !util.PathExists(videoFilePath) {
			if err := util.DownloadFileMultiThreading(url, videoFilePath, 100*1024*1024, 8, nil); err != nil {
				log.Errorf("failed to download video file, err: %+v", err)
				return EmptyText()
			}
		}
	} else {
		videoFilePath = url
	}

	videoBytes, err := ioutil.ReadFile(videoFilePath)
	if err != nil {
		log.Errorf("failed to open local video file, %+v", err)
		return EmptyText()
	}
	elem.UploadingVideo = bytes.NewReader(videoBytes)
	elem.Url = url           // 仅用于发送日志展示
	elem.CoverUrl = coverUrl // 仅用于发送日志展示
	return elem
}

func ProtoMusicToMiraiMusic(_ *client.QQClient, data map[string]string) (m message.IMessageElement) {
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
}
