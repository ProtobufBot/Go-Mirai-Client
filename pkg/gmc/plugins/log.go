package plugins

import (
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/bot"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/plugin"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	log "github.com/sirupsen/logrus"
)

func LogPrivateMessage(cli *client.QQClient, event *message.PrivateMessage) int32 {
	log.Infof("Bot(%+v) Private(%+v) -> %+v\n", cli.Uin, event.Sender.Uin, bot.MiraiMsgToRawMsg(cli, event.Elements))
	return plugin.MessageIgnore
}

func LogGroupMessage(cli *client.QQClient, event *message.GroupMessage) int32 {
	cli.MarkGroupMessageReaded(event.GroupCode, int64(event.Id)) // 标记为已读，可能可以减少风控
	log.Infof("Bot(%+v) Group(%+v) -> %+v\n", cli.Uin, event.GroupCode, bot.MiraiMsgToRawMsg(cli, event.Elements))
	return plugin.MessageIgnore
}
