package plugins

import (
	"time"

	"github.com/2mf8/Go-Lagrange-Client/pkg/bot"
	"github.com/2mf8/Go-Lagrange-Client/pkg/cache"
	"github.com/2mf8/Go-Lagrange-Client/pkg/plugin"
	"github.com/2mf8/Go-Lagrange-Client/proto_gen/onebot"

	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/event"
	"github.com/LagrangeDev/LagrangeGo/message"
)

func ReportPrivateMessage(cli *client.QQClient, event *message.PrivateMessage) int32 {
	cache.PrivateMessageLru.Add(event.Id, event)
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TPrivateMessageEvent,
	}
	eventProto.Data = &onebot.Frame_PrivateMessageEvent{
		PrivateMessageEvent: &onebot.PrivateMessageEvent{
			Time:        time.Now().Unix(),
			SelfId:      int64(cli.Uin),
			PostType:    "message",
			MessageType: "private",
			SubType:     "normal",
			MessageId:   event.Id,
			MessageReceipt: &onebot.MessageReceipt{
				SenderId: int64(event.Sender.Uin),
				Time:     time.Now().Unix(),
				Seqs:     []int32{event.Id},
			},
			UserId:     int64(event.Sender.Uin),
			Message:    bot.MiraiMsgToProtoMsg(cli, event.Elements),
			RawMessage: bot.MiraiMsgToRawMsg(cli, event.Elements),
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
		MessageReceipt: &onebot.MessageReceipt{
			SenderId: int64(event.Sender.Uin),
			Time:     time.Now().Unix(),
			Seqs:     []int32{event.Id},
			GroupId:  int64(event.GroupCode),
		},
		GroupId:    int64(event.GroupCode),
		UserId:     int64(event.Sender.Uin),
		Message:    bot.MiraiMsgToProtoMsg(cli, event.Elements),
		RawMessage: bot.MiraiMsgToRawMsg(cli, event.Elements),
		Sender: &onebot.GroupMessageEvent_Sender{
			UserId: int64(event.Sender.Uin),
		},
	}

	eventProto.Data = &onebot.Frame_GroupMessageEvent{
		GroupMessageEvent: groupMessageEvent,
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportMemberJoin(cli *client.QQClient, event *event.GroupMemberIncrease) int32 {
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TGroupIncreaseNoticeEvent,
	}
	eventProto.Data = &onebot.Frame_GroupIncreaseNoticeEvent{
		GroupIncreaseNoticeEvent: &onebot.GroupIncreaseNoticeEvent{
			Time:       time.Now().Unix(),
			SelfId:     int64(cli.Uin),
			PostType:   "message",
			NoticeType: "group_increase",
			SubType:    "approve",
			GroupId:    int64(event.GroupUin),
			UserId:     0,
			OperatorId: 0,
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

/*func ReportMemberLeave(cli *client.QQClient, event *event.MemberLeaveGroupEvent) int32 {
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

func ReportJoinGroup(cli *client.QQClient, event *event.GroupInfo) int32 {
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

func ReportLeaveGroup(cli *client.QQClient, event *event.GroupLeaveEvent) int32 {
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

func ReportGroupMute(cli *client.QQClient, event *event.GroupMuteEvent) int32 {
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TGroupBanNoticeEvent,
	}
	eventProto.Data = &onebot.Frame_GroupBanNoticeEvent{
		GroupBanNoticeEvent: &onebot.GroupBanNoticeEvent{
			Time:       time.Now().Unix(),
			SelfId:     cli.Uin,
			PostType:   "notice",
			NoticeType: "group_ban",
			SubType: func() string {
				if event.Time == 0 {
					return "lift_ban"
				}
				return "ban"
			}(),
			GroupId:    event.GroupCode,
			OperatorId: event.OperatorUin,
			UserId:     event.TargetUin,
			Duration:   int64(event.Time),
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportNewFriendRequest(cli *client.QQClient, event *event.NewFriendRequest) int32 {
	flag := strconv.FormatInt(event.RequestId, 10)
	cache.FriendRequestLru.Add(flag, event)
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TFriendRequestEvent,
	}
	eventProto.Data = &onebot.Frame_FriendRequestEvent{
		FriendRequestEvent: &onebot.FriendRequestEvent{
			Time:        time.Now().Unix(),
			SelfId:      cli.Uin,
			PostType:    "request",
			RequestType: "friend",
			UserId:      event.RequesterUin,
			Comment:     event.Message,
			Flag:        flag,
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportUserJoinGroupRequest(cli *client.QQClient, event *event.UserJoinGroupRequest) int32 {
	flag := strconv.FormatInt(event.RequestId, 10)
	cache.GroupRequestLru.Add(flag, event)
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TGroupRequestEvent,
	}
	eventProto.Data = &onebot.Frame_GroupRequestEvent{
		GroupRequestEvent: &onebot.GroupRequestEvent{
			Time:          time.Now().Unix(),
			SelfId:        cli.Uin,
			PostType:      "request",
			RequestType:   "group",
			SubType:       "add",
			GroupId:       event.GroupCode,
			UserId:        event.RequesterUin,
			Comment:       event.Message,
			Flag:          flag,
			RequestId:     event.RequestId,
			UserNick:      event.RequesterNick,
			ActionUinNick: event.ActionUinNick,
			ActionUin:     event.ActionUin,
			Check:         event.Checked,
			Suspicious:    event.Suspicious,
			Extra: map[string]string{
				"actor_id": strconv.FormatInt(event.Actor, 10),
			},
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportGroupInvitedRequest(cli *client.QQClient, event *event.GroupInvitedRequest) int32 {
	flag := strconv.FormatInt(event.RequestId, 10)
	cache.GroupInvitedRequestLru.Add(flag, event)
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TGroupRequestEvent,
	}
	eventProto.Data = &onebot.Frame_GroupRequestEvent{
		GroupRequestEvent: &onebot.GroupRequestEvent{
			Time:        time.Now().Unix(),
			SelfId:      cli.Uin,
			PostType:    "request",
			RequestType: "group",
			SubType:     "invite",
			GroupId:     event.GroupCode,
			UserId:      event.InvitorUin,
			Comment:     "",
			Flag:        flag,
			Extra: map[string]string{
				"actor_id": strconv.FormatInt(event.Actor, 10),
			},
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportGroupMessageRecalled(cli *client.QQClient, event *event.GroupMessageRecalledEvent) int32 {
	if event.AuthorUin == event.OperatorUin {
		log.Infof("群 %v 内 %v 撤回了一条消息, 消息Id为 %v", event.GroupCode, event.AuthorUin, event.MessageId)
	} else {
		log.Infof("群 %v 内 %v 撤回了 %v 的一条消息, 消息Id为 %v", event.GroupCode, event.OperatorUin, event.AuthorUin, event.MessageId)
	}
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TGroupRecallNoticeEvent,
	}
	eventProto.Data = &onebot.Frame_GroupRecallNoticeEvent{
		GroupRecallNoticeEvent: &onebot.GroupRecallNoticeEvent{
			Time:       time.Now().Unix(),
			SelfId:     cli.Uin,
			PostType:   "notice",
			NoticeType: "group_recall",
			GroupId:    event.GroupCode,
			UserId:     event.AuthorUin,
			OperatorId: event.OperatorUin,
			MessageId:  event.MessageId,
			MessageReceipt: &onebot.MessageReceipt{
				SenderId: event.AuthorUin,
				Time:     time.Now().Unix(),
				Seqs:     []int32{event.MessageId},
				GroupId:  event.GroupCode,
			},
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportFriendMessageRecalled(cli *client.QQClient, event *event.FriendMessageRecalledEvent) int32 {
	log.Infof("好友 %v 撤回了一条消息, 消息Id为 %v", event.FriendUin, event.MessageId)
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TFriendRecallNoticeEvent,
	}
	eventProto.Data = &onebot.Frame_FriendRecallNoticeEvent{
		FriendRecallNoticeEvent: &onebot.FriendRecallNoticeEvent{
			Time:       time.Now().Unix(),
			SelfId:     cli.Uin,
			PostType:   "notice",
			NoticeType: "friend_recall",
			UserId:     event.FriendUin,
			MessageId:  event.MessageId,
			MessageReceipt: &onebot.MessageReceipt{
				SenderId: event.FriendUin,
				Time:     time.Now().Unix(),
				Seqs:     []int32{event.MessageId},
			},
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportNewFriendAdded(cli *client.QQClient, event *event.NewFriendEvent) int32 {
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TFriendAddNoticeEvent,
	}
	eventProto.Data = &onebot.Frame_FriendAddNoticeEvent{
		FriendAddNoticeEvent: &onebot.FriendAddNoticeEvent{
			Time:       time.Now().Unix(),
			SelfId:     cli.Uin,
			PostType:   "notice",
			NoticeType: "friend_add",
			UserId:     event.Friend.Uin,
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

// 暂时先放在私聊里面吧，onebot协议里面没这个
func ReportOfflineFile(cli *client.QQClient, event *event.OfflineFileEvent) int32 {
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TPrivateMessageEvent,
	}
	eventProto.Data = &onebot.Frame_PrivateMessageEvent{
		PrivateMessageEvent: &onebot.PrivateMessageEvent{
			Time:        time.Now().Unix(),
			SelfId:      cli.Uin,
			PostType:    "message",
			MessageType: "private",
			SubType:     "normal",
			MessageId:   0,
			MessageReceipt: &onebot.MessageReceipt{
				SenderId: event.Sender,
				Time:     time.Now().Unix(),
				Seqs:     []int32{0},
			},
			UserId: event.Sender,
			Message: []*onebot.Message{
				{
					Type: "file",
					Data: map[string]string{
						"url":  event.DownloadUrl,
						"name": event.FileName,
						"size": strconv.FormatInt(event.FileSize, 10),
					},
				},
			},
			RawMessage: fmt.Sprintf(`<file url="%s" name="%s" size="%d"/>`, html.EscapeString(event.DownloadUrl), html.EscapeString(event.FileName), event.FileSize),
			Sender: &onebot.PrivateMessageEvent_Sender{
				UserId: event.Sender,
			},
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportGroupNotify(cli *client.QQClient, event client.INotifyEvent) int32 {
	group := cli.FindGroup(event.From())
	switch notify := event.(type) {
	case *client.GroupPokeNotifyEvent:
		sender := group.FindMember(notify.Sender)
		receiver := group.FindMember(notify.Receiver)
		if sender == receiver {
			log.Infof("群 %v(%v) 内 %v(%v) 戳了戳自己", group.Code, group.Name, sender.Uin, sender.Nickname)
		} else {
			log.Infof("群 %v(%v) 内 %v(%v) 戳了戳 %v(%v)", group.Code, group.Name, sender.Uin, sender.Nickname, receiver.Uin, receiver.Nickname)
		}
		eventProto := &onebot.Frame{
			FrameType: onebot.Frame_TGroupNotifyEvent,
		}
		eventProto.Data = &onebot.Frame_GroupNotifyEvent{
			GroupNotifyEvent: &onebot.GroupNotifyEvent{
				Time:       time.Now().Unix(),
				SelfId:     cli.Uin,
				PostType:   "notice",
				NoticeType: "group_poke",
				GroupId:    group.Code,
				GroupName:  group.Name,
				Sender:     sender.Uin,
				SenderCard: sender.Nickname,
				TargetId:   receiver.Uin,
				TargetCard: receiver.Nickname,
			},
		}
		bot.HandleEventFrame(cli, eventProto)
		return plugin.MessageIgnore
	case *client.GroupRedBagLuckyKingNotifyEvent:
		sender := group.FindMember(notify.Sender)
		luckyKing := group.FindMember(notify.LuckyKing)
		log.Infof("群 %v(%v) 内 %v(%v) 的红包被抢完, %v(%v) 是运气王", group.Code, group.Name, sender.Uin, sender.Nickname, luckyKing.Uin, luckyKing.Nickname)
		eventProto := &onebot.Frame{
			FrameType: onebot.Frame_TGroupNotifyEvent,
		}
		eventProto.Data = &onebot.Frame_GroupNotifyEvent{
			GroupNotifyEvent: &onebot.GroupNotifyEvent{
				Time:       time.Now().Unix(),
				SelfId:     cli.Uin,
				PostType:   "notice",
				NoticeType: "group_red_bag_lucky_king",
				GroupId:    group.Code,
				GroupName:  group.Name,
				Sender:     sender.Uin,
				SenderCard: sender.Nickname,
				TargetId:   luckyKing.Uin,
				TargetCard: luckyKing.Nickname,
			},
		}
		bot.HandleEventFrame(cli, eventProto)
		return plugin.MessageIgnore
	case *client.MemberHonorChangedNotifyEvent:
		log.Info(notify.Content())
		eventProto := &onebot.Frame{
			FrameType: onebot.Frame_TGroupNotifyEvent,
		}
		eventProto.Data = &onebot.Frame_GroupNotifyEvent{
			GroupNotifyEvent: &onebot.GroupNotifyEvent{
				Time:       time.Now().Unix(),
				SelfId:     cli.Uin,
				PostType:   "notice",
				NoticeType: "member_honor_change",
				GroupId:    group.Code,
				GroupName:  group.Name,
				TargetId:   notify.Uin,
				TargetCard: notify.Nick,
				Honor: func() string {
					switch notify.Honor {
					case client.Talkative:
						return "talkative"
					case client.Performer:
						return "performer"
					case client.Emotion:
						return "emotion"
					case client.Legend:
						return "legend"
					case client.StrongNewbie:
						return "strong_newbie"
					default:
						return "ERROR"
					}
				}(),
			},
		}
		bot.HandleEventFrame(cli, eventProto)
		return plugin.MessageIgnore
	}
	return plugin.MessageIgnore
}*/
