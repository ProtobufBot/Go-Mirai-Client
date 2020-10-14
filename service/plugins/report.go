package plugins

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/plugin"
	"github.com/ProtobufBot/Go-Mirai-Client/proto_gen/onebot"
	"github.com/ProtobufBot/Go-Mirai-Client/service/bot"
	"time"
)

func ReportPrivateMessage(cli *client.QQClient, event *message.PrivateMessage) int32 {
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TPrivateMessageEvent,
	}
	eventProto.Data = &onebot.Frame_PrivateMessageEvent{
		PrivateMessageEvent: &onebot.PrivateMessageEvent{
			PostType:    "message",
			MessageType: "private",
			SubType:     "normal",
			MessageId:   event.Id,
			UserId:      event.Sender.Uin,
			Message:     bot.MiraiMsgToProtoMsg(event.Elements),
			RawMessage:  bot.MiraiMsgToRawMsg(event.Elements),
			Sender: &onebot.PrivateMessageEvent_Sender{
				UserId:   event.Sender.Uin,
				Nickname: event.Sender.Nickname,
			},
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportGroupMessage(cli *client.QQClient, event *message.GroupMessage) int32 {
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TGroupMessageEvent,
	}
	eventProto.Data = &onebot.Frame_GroupMessageEvent{
		GroupMessageEvent: &onebot.GroupMessageEvent{
			Time:        time.Now().Unix(),
			SelfId:      cli.Uin,
			PostType:    "message",
			MessageType: "group",
			SubType:     "normal",
			MessageId:   event.Id,
			GroupId:     event.GroupCode,
			UserId:      event.Sender.Uin,
			Message:     bot.MiraiMsgToProtoMsg(event.Elements),
			RawMessage:  bot.MiraiMsgToRawMsg(event.Elements),
			Sender: &onebot.GroupMessageEvent_Sender{
				UserId:   event.Sender.Uin,
				Nickname: event.Sender.Nickname,
			},
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportMemberJoin(cli *client.QQClient, event *client.MemberJoinGroupEvent) int32 {
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TGroupIncreaseNoticeEvent,
	}
	eventProto.Data = &onebot.Frame_GroupIncreaseNoticeEvent{
		GroupIncreaseNoticeEvent: &onebot.GroupIncreaseNoticeEvent{
			Time:       time.Now().Unix(),
			SelfId:     cli.Uin,
			PostType:   "message",
			NoticeType: "group_increase",
			SubType:    "approve",
			GroupId:    event.Group.Code,
			UserId:     event.Member.Uin,
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportMemberLeave(cli *client.QQClient, event *client.MemberLeaveGroupEvent) int32 {
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TGroupDecreaseNoticeEvent,
	}
	subType := "leave"
	var operatorId int64 = 0
	if event.Operator != nil {
		subType = "kick"
		operatorId = event.Operator.Uin
	}

	eventProto.Data = &onebot.Frame_GroupDecreaseNoticeEvent{
		GroupDecreaseNoticeEvent: &onebot.GroupDecreaseNoticeEvent{
			Time:       time.Now().Unix(),
			SelfId:     cli.Uin,
			PostType:   "message",
			NoticeType: "group_decrease",
			SubType:    subType,
			GroupId:    event.Group.Code,
			UserId:     event.Member.Uin,
			OperatorId: operatorId,
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}
