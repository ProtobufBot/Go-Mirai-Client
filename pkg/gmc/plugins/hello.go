package plugins

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/2mf8/Go-Lagrange-Client/pkg/bot"
	"github.com/2mf8/Go-Lagrange-Client/pkg/plugin"
	"github.com/2mf8/LagrangeGo/client"
	"github.com/2mf8/LagrangeGo/message"
	log "github.com/sirupsen/logrus"
)

func HelloPrivateMessage(cli *client.QQClient, event *message.PrivateMessage) int32 {
	b, _ := json.Marshal(event)
	log.Warn(string(b))
	if bot.MiraiMsgToRawMsg(cli, event.Elements) != "hi" {
		return plugin.MessageIgnore
	}
	elem := &message.SendingMessage{
		Elements: []message.IMessageElement{
			&message.TextElement{Content: "hello"},
		},
	}
	cli.SendPrivateMessage(event.Sender.Uin, elem.Elements)
	return plugin.MessageIgnore
}

func HelloGroupMessage(cli *client.QQClient, event *message.GroupMessage) int32 {
	if bot.MiraiMsgToRawMsg(cli, event.Elements) != "hi" {
		return plugin.MessageIgnore
	}
	resp, err := http.Get("https://2mf8.cn/logo.png")
	defer resp.Body.Close()
	if err != nil {
		return plugin.MessageIgnore
	}
	imo, err := io.ReadAll(resp.Body)
	if err != nil {
		return plugin.MessageIgnore
	}
	ir, err := cli.ImageUploadGroup(event.GroupCode, &message.GroupImageElement{
		Stream: imo,
	})
	if err != nil {
		return plugin.MessageIgnore
	}
	elem := &message.SendingMessage{}
	elem.Elements = append(elem.Elements, ir)
	elem.Elements = append(elem.Elements, &message.TextElement{
		Content: "测试成功",
	})
	r, e := cli.SendGroupMessage(event.GroupCode, elem.Elements)
	log.Warn(r, e)
	return plugin.MessageIgnore
}
