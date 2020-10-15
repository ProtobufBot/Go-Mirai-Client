package plugins

import (
	"time"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/plugin"
	"github.com/ProtobufBot/Go-Mirai-Client/proto_gen/onebot"
	"github.com/ProtobufBot/Go-Mirai-Client/service/bot"
	"github.com/ProtobufBot/Go-Mirai-Client/service/cache"
)

func ReportPrivateMessage(cli *client.QQClient, event *message.PrivateMessage) int32 {
	messageId := cache.NextGlobalSeq()
	cache.PrivateMessageLru.Add(messageId, event)
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TPrivateMessageEvent,
	}
	eventProto.Data = &onebot.Frame_PrivateMessageEvent{
		PrivateMessageEvent: &onebot.PrivateMessageEvent{
			PostType:    "message",
			MessageType: "private",
			SubType:     "normal",
			MessageId:   messageId,
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
	messageId := cache.NextGlobalSeq()
	cache.GroupMessageLru.Add(messageId, event)
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
			MessageId:   messageId,
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
			OperatorId: 0,
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

func ReportJoinGroup(cli *client.QQClient, event *client.GroupInfo) int32 {
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
			GroupId:    event.Code,
			UserId:     cli.Uin,
			OperatorId: 0,
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportLeaveGroup(cli *client.QQClient, event *client.GroupLeaveEvent) int32 {
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
			UserId:     cli.Uin,
			OperatorId: operatorId,
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}
