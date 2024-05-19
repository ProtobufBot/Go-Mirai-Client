package plugins

import (
	"strconv"
	"time"

	"github.com/2mf8/Go-Lagrange-Client/pkg/bot"
	"github.com/2mf8/Go-Lagrange-Client/pkg/cache"
	"github.com/2mf8/Go-Lagrange-Client/pkg/plugin"
	"github.com/2mf8/Go-Lagrange-Client/proto_gen/onebot"
	log "github.com/sirupsen/logrus"

	"github.com/2mf8/LagrangeGo/client"
	"github.com/2mf8/LagrangeGo/client/event"
	"github.com/2mf8/LagrangeGo/message"
)

func ReportPrivateMessage(cli *client.QQClient, event *message.PrivateMessage) int32 {
	cache.PrivateMessageLru.Add(event.Id, event)
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TPrivateMessageEvent,
	}
	eventProto.PbData = &onebot.Frame_PrivateMessageEvent{
		PrivateMessageEvent: &onebot.PrivateMessageEvent{
			Time:        time.Now().Unix(),
			SelfId:      int64(cli.Uin),
			PostType:    "message",
			MessageType: "private",
			SubType:     "normal",
			MessageId:   event.Id,
			UserId:      int64(event.Sender.Uin),
			Message:     bot.MiraiMsgToProtoMsg(cli, event.Elements),
			RawMessage:  bot.MiraiMsgToRawMsg(cli, event.Elements),
			Sender: &onebot.PrivateMessageEvent_Sender{
				UserId:   int64(event.Sender.Uin),
				Nickname: event.Sender.Nickname,
			},
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportGroupMessage(cli *client.QQClient, event *message.GroupMessage) int32 {
	cache.GroupMessageLru.Add(event.Id, event)
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TGroupMessageEvent,
	}
	groupMessageEvent := &onebot.GroupMessageEvent{
		Time:        time.Now().Unix(),
		SelfId:      int64(cli.Uin),
		PostType:    "message",
		MessageType: "group",
		SubType:     "normal",
		MessageId:   event.Id,
		GroupId:     int64(event.GroupCode),
		UserId:      int64(event.Sender.Uin),
		Message:     bot.MiraiMsgToProtoMsg(cli, event.Elements),
		RawMessage:  bot.MiraiMsgToRawMsg(cli, event.Elements),
		Sender: &onebot.GroupMessageEvent_Sender{
			UserId: int64(event.Sender.Uin),
			Nickname: event.Sender.Nickname,
			Card: event.Sender.CardName,
		},
	}

	eventProto.PbData = &onebot.Frame_GroupMessageEvent{
		GroupMessageEvent: groupMessageEvent,
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportMemberJoin(cli *client.QQClient, event *event.GroupMemberIncrease) int32 {
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TGroupIncreaseNoticeEvent,
	}
	eventProto.PbData = &onebot.Frame_GroupIncreaseNoticeEvent{
		GroupIncreaseNoticeEvent: &onebot.GroupIncreaseNoticeEvent{
			Time:       time.Now().Unix(),
			SelfId:     int64(cli.Uin),
			PostType:   "message",
			NoticeType: "group_increase",
			SubType:    "approve",
			GroupId:    int64(event.GroupUin),
			UserId:     0,
			OperatorId: 0,
			MemberUid:  event.MemberUid,
			InvitorUid: event.InvitorUid,
			JoinType:   event.JoinType,
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportMemberLeave(cli *client.QQClient, event *event.GroupMemberDecrease) int32 {
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TGroupDecreaseNoticeEvent,
	}
	subType := "leave"
	var operatorUid string = ""
	if event.IsKicked() {
		subType = "kick"
		operatorUid = event.OperatorUid
	}

	eventProto.PbData = &onebot.Frame_GroupDecreaseNoticeEvent{
		GroupDecreaseNoticeEvent: &onebot.GroupDecreaseNoticeEvent{
			Time:        time.Now().Unix(),
			SelfId:      int64(cli.Uin),
			PostType:    "message",
			NoticeType:  "group_decrease",
			SubType:     subType,
			GroupId:     int64(event.GroupUin),
			MemberUid:   event.MemberUid,
			OperatorUid: operatorUid,
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportJoinGroup(cli *client.QQClient, event *event.GroupMemberIncrease) int32 {
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TGroupIncreaseNoticeEvent,
	}
	eventProto.PbData = &onebot.Frame_GroupIncreaseNoticeEvent{
		GroupIncreaseNoticeEvent: &onebot.GroupIncreaseNoticeEvent{
			Time:       time.Now().Unix(),
			SelfId:     int64(cli.Uin),
			PostType:   "message",
			NoticeType: "group_increase",
			SubType:    "approve",
			GroupId:    int64(event.GroupUin),
			UserId:     int64(cli.Uin),
			OperatorId: 0,
			MemberUid:  event.MemberUid,
			JoinType:   event.JoinType,
			InvitorUid: event.InvitorUid,
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportGroupMute(cli *client.QQClient, event *event.GroupMute) int32 {
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TGroupBanNoticeEvent,
	}
	eventProto.PbData = &onebot.Frame_GroupBanNoticeEvent{
		GroupBanNoticeEvent: &onebot.GroupBanNoticeEvent{
			Time:       time.Now().Unix(),
			SelfId:     int64(cli.Uin),
			PostType:   "notice",
			NoticeType: "group_ban",
			SubType: func() string {
				if event.Duration == 0 {
					return "lift_ban"
				}
				return "ban"
			}(),
			GroupId:     int64(event.GroupUin),
			OperatorUid: event.OperatorUid,
			TargetUid:   event.TargetUid,
			Duration:    int64(event.Duration),
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportNewFriendRequest(cli *client.QQClient, event *event.FriendRequest) int32 {
	flag := strconv.FormatInt(int64(event.SourceUin), 10)
	cache.FriendRequestLru.Add(flag, event)
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TFriendRequestEvent,
	}
	eventProto.PbData = &onebot.Frame_FriendRequestEvent{
		FriendRequestEvent: &onebot.FriendRequestEvent{
			Time:        time.Now().Unix(),
			SelfId:      int64(cli.Uin),
			PostType:    "request",
			RequestType: "friend",
			Flag:        flag,
			SourceUid:   event.SourceUid,
			Msg:         event.Msg,
			Source:      event.Source,
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportUserJoinGroupRequest(cli *client.QQClient, event *event.GroupMemberJoinRequest) int32 {
	flag := strconv.FormatInt(int64(event.GroupUin), 10)
	cache.GroupRequestLru.Add(flag, event)
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TGroupRequestEvent,
	}
	eventProto.PbData = &onebot.Frame_GroupRequestEvent{
		GroupRequestEvent: &onebot.GroupRequestEvent{
			Time:        time.Now().Unix(),
			SelfId:      int64(cli.Uin),
			PostType:    "request",
			RequestType: "group",
			SubType:     "add",
			GroupId:     int64(event.GroupUin),
			Flag:        flag,
			TargetUid:   event.TargetUid,
			InvitorUid:  event.InvitorUid,
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportGroupInvitedRequest(cli *client.QQClient, event *event.GroupInvite) int32 {
	flag := strconv.FormatInt(int64(event.GroupUin), 10)
	cache.GroupInvitedRequestLru.Add(flag, event)
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TGroupRequestEvent,
	}
	eventProto.PbData = &onebot.Frame_GroupRequestEvent{
		GroupRequestEvent: &onebot.GroupRequestEvent{
			Time:        time.Now().Unix(),
			SelfId:      int64(cli.Uin),
			PostType:    "request",
			RequestType: "group",
			SubType:     "invite",
			GroupId:     int64(event.GroupUin),
			InvitorUid:  event.InvitorUid,
			Comment:     "",
			Flag:        flag,
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportGroupMessageRecalled(cli *client.QQClient, event *event.GroupRecall) int32 {
	if event.AuthorUid == event.OperatorUid {
		log.Infof("群 %v 内 %s 撤回了一条消息, 消息Id为 %v", event.GroupUin, event.AuthorUid, event.Sequence)
	} else {
		log.Infof("群 %v 内 %s 撤回了 %s 的一条消息, 消息Id为 %v", event.GroupUin, event.OperatorUid, event.AuthorUid, event.Sequence)
	}
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TGroupRecallNoticeEvent,
	}
	eventProto.PbData = &onebot.Frame_GroupRecallNoticeEvent{
		GroupRecallNoticeEvent: &onebot.GroupRecallNoticeEvent{
			Time:        time.Now().Unix(),
			SelfId:      int64(cli.Uin),
			PostType:    "notice",
			NoticeType:  "group_recall",
			GroupId:     int64(event.GroupUin),
			AuthorUid:   event.AuthorUid,
			OperatorUid: event.OperatorUid,
			Sequence:    event.Sequence,
			Random:      event.Random,
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportFriendMessageRecalled(cli *client.QQClient, event *event.FriendRecall) int32 {
	log.Infof("好友 %s 撤回了一条消息, 消息Id为 %v", event.FromUid, event.Sequence)
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TFriendRecallNoticeEvent,
	}
	eventProto.PbData = &onebot.Frame_FriendRecallNoticeEvent{
		FriendRecallNoticeEvent: &onebot.FriendRecallNoticeEvent{
			Time:       time.Now().Unix(),
			SelfId:     int64(cli.Uin),
			PostType:   "notice",
			NoticeType: "friend_recall",
			FromUid:    event.FromUid,
			MessageId:  int32(event.Sequence),
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportNewFriendAdded(cli *client.QQClient, event *event.FriendRequest) int32 {
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TFriendAddNoticeEvent,
	}
	eventProto.PbData = &onebot.Frame_FriendAddNoticeEvent{
		FriendAddNoticeEvent: &onebot.FriendAddNoticeEvent{
			Time:       time.Now().Unix(),
			SelfId:     int64(cli.Uin),
			PostType:   "notice",
			NoticeType: "friend_add",
			UserId:     int64(event.SourceUin),
			SourceUin:  event.SourceUin,
			SourceUid:  event.SourceUid,
			Source:     event.Source,
			Msg:        event.Msg,
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}
