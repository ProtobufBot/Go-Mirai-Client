package plugins

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/message"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/bot"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/plugin"
)

func HelloPrivateMessage(cli *client.QQClient, event *message.PrivateMessage) int32 {
	b, _:= json.Marshal(event)
	log.Warn(string(b))
	if bot.MiraiMsgToRawMsg(cli,event.Elements) != "hi" {
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