package plugins

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/plugin"
	"github.com/ProtobufBot/Go-Mirai-Client/service/bot"
)

func HelloPrivateMessage(cli *client.QQClient, event *message.PrivateMessage) int32 {
	if event.Sender.Uin != 875543533 {
		return plugin.MessageIgnore
	}
	if bot.MiraiMsgToRawMsg(event.Elements) != "hi" {
		return plugin.MessageIgnore
	}
	cli.SendPrivateMessage(event.Sender.Uin, &message.SendingMessage{
		Elements: []message.IMessageElement{
			&message.TextElement{Content: "hello"},
		},
	})
	return plugin.MessageIgnore
}
