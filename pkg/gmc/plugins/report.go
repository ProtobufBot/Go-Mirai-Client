package plugins

import (
	"fmt"
	"html"
	"strconv"
	"time"

	"github.com/ProtobufBot/Go-Mirai-Client/pkg/bot"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/cache"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/plugin"
	"github.com/ProtobufBot/Go-Mirai-Client/proto_gen/onebot"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	log "github.com/sirupsen/logrus"
)

func ReportPrivateMessage(cli *client.QQClient, event *message.PrivateMessage) int32 {
	cache.PrivateMessageLru.Add(event.Id, event)
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
			MessageId:   event.Id,
			MessageReceipt: &onebot.MessageReceipt{
				SenderId: event.Sender.Uin,
				Time:     time.Now().Unix(),
				Seqs:     []int32{event.Id},
			},
			UserId:     event.Sender.Uin,
			Message:    bot.MiraiMsgToProtoMsg(cli, event.Elements),
			RawMessage: bot.MiraiMsgToRawMsg(cli, event.Elements),
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
	if CheckGroupFile(cli, event) { // 检查是否有群文件element，如果有，执行GroupUploadNotice
		return plugin.MessageIgnore
	}
	cache.GroupMessageLru.Add(event.Id, event)
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TGroupMessageEvent,
	}
	groupMessageEvent := &onebot.GroupMessageEvent{
		Time:        time.Now().Unix(),
		SelfId:      cli.Uin,
		PostType:    "message",
		MessageType: "group",
		SubType:     "normal",
		MessageId:   event.Id,
		MessageReceipt: &onebot.MessageReceipt{
			SenderId: event.Sender.Uin,
			Time:     time.Now().Unix(),
			Seqs:     []int32{event.Id},
			GroupId:  event.GroupCode,
		},
		GroupId:    event.GroupCode,
		UserId:     event.Sender.Uin,
		Message:    bot.MiraiMsgToProtoMsg(cli, event.Elements),
		RawMessage: bot.MiraiMsgToRawMsg(cli, event.Elements),
		Sender: &onebot.GroupMessageEvent_Sender{
			UserId: event.Sender.Uin,
		},
	}
	if event.Sender.IsAnonymous() {
		groupMessageEvent.SubType = "anonymous"
		groupMessageEvent.Anonymous = &onebot.GroupMessageEvent_Anonymous{
			Name: event.Sender.Nickname,
		}
	} else {
		group := cli.FindGroup(event.GroupCode)
		if group == nil {
			err := cli.ReloadGroupList()
			group = cli.FindGroup(event.GroupCode)
			if err != nil || group == nil {
				log.Warnf("failed to find group: %+v, err: %+v", event.GroupCode, err)
				return plugin.MessageIgnore
			}
		}
		member := group.FindMember(event.Sender.Uin)
		if member != nil {
			groupMessageEvent.Sender.Role = func() string {
				switch member.Permission {
				case client.Owner:
					return "owner"
				case client.Administrator:
					return "admin"
				default:
					return "member"
				}
			}()
			groupMessageEvent.Sender.Nickname = member.Nickname
			groupMessageEvent.Sender.Title = member.SpecialTitle
			groupMessageEvent.Sender.Card = member.CardName
		}
	}

	eventProto.Data = &onebot.Frame_GroupMessageEvent{
		GroupMessageEvent: groupMessageEvent,
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportChannelMessage(cli *client.QQClient, event *message.GuildChannelMessage) int32 {
	cache.ChannelMessageLru.Add(event.Id, event)
	v, ok := cache.GetGuildAdminTimeLru.Get(event.GuildId)
	gtime, _ := strconv.ParseInt(fmt.Sprintf("%v", v), 10, 64)
	if !ok || time.Now().Unix() > gtime {
		log.Println("缓存频道管理员更新")
		cache.GetGuildAdminTimeLru.Remove(event.GuildId)
		users, _ := cli.GuildService.FetchGuildMemberListWithRole(event.GuildId, event.ChannelId, 0, 2, "")
		for _, v := range users.Members {
			// log.Println(v.Nickname, v.Role, v.RoleName, v.TinyId, v.Title, v.LastSpeakTime)
			if v.Role == 2 {
				cache.GuildAdminLru.Remove(v.TinyId)
				cache.GuildAdminLru.Add(v.TinyId, v)
			} else {
				break
			}
		}
		cache.GetGuildAdminTimeLru.Add(event.GuildId, time.Now().Add(time.Minute*10).Unix())
	}

	role, _ := cache.GuildAdminLru.Get(event.Sender.TinyId)
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TChannelMessageEvent,
	}
	channelMessageEvent := &onebot.ChannelMessageEvent{
		Id:          event.Id,
		InternalId:  event.InternalId,
		GuildId:     event.GuildId,
		ChannelId:   event.ChannelId,
		Time:        event.Time,
		SelfId:      int64(cli.GuildService.TinyId),
		PostType:    "message",
		MessageType: "channel",
		SubType:     "normal",
		Message:     bot.MiraiMsgToProtoMsg(cli, event.Elements),
		RawMessage:  bot.MiraiMsgToRawMsg(cli, event.Elements),
		Sender: &onebot.ChannelMessageEvent_Sender{
			TinyId:   event.Sender.TinyId,
			Nickname: event.Sender.Nickname,
		},
		MessageId: event.Id,
	}
	if role != nil {
		channelMessageEvent.Sender.Roles = append(channelMessageEvent.Sender.Roles, 2)
		channelMessageEvent.Sender.RoleNames = append(channelMessageEvent.Sender.RoleNames, "管理员")
	}
	eventProto.Data = &onebot.Frame_ChannelMessageEvent{
		ChannelMessageEvent: channelMessageEvent,
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func CheckGroupFile(cli *client.QQClient, event *message.GroupMessage) bool {
	for _, elem := range event.Elements {
		if file, ok := elem.(*message.GroupFileElement); ok {
			eventProto := &onebot.Frame{
				FrameType: onebot.Frame_TGroupUploadNoticeEvent,
			}
			groupUploadNoticeEvent := &onebot.GroupUploadNoticeEvent{
				Time:       time.Now().Unix(),
				SelfId:     cli.Uin,
				PostType:   "notice",
				NoticeType: "group_upload",
				GroupId:    event.GroupCode,
				UserId:     event.Sender.Uin,
				File: &onebot.GroupUploadNoticeEvent_File{
					Id:    file.Path,
					Name:  file.Name,
					Busid: int64(file.Busid),
					Size:  file.Size,
					Url:   cli.GetGroupFileUrl(event.GroupCode, file.Path, file.Busid),
				},
			}
			eventProto.Data = &onebot.Frame_GroupUploadNoticeEvent{
				GroupUploadNoticeEvent: groupUploadNoticeEvent,
			}
			bot.HandleEventFrame(cli, eventProto)
			return true
		}
	}
	return false
}

func ReportTempMessage(cli *client.QQClient, event *client.TempMessageEvent) int32 {
	// TODO 撤回？
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TGroupTempMessageEvent,
	}
	eventProto.Data = &onebot.Frame_GroupTempMessageEvent{
		GroupTempMessageEvent: &onebot.GroupTempMessageEvent{
			Time:        time.Now().Unix(),
			SelfId:      cli.Uin,
			PostType:    "message",
			MessageType: "group_temp",
			MessageId:   event.Message.Id,
			MessageReceipt: &onebot.MessageReceipt{
				SenderId: event.Message.Sender.Uin,
				Time:     time.Now().Unix(),
				Seqs:     []int32{event.Message.Id},
				GroupId:  event.Message.GroupCode,
			},
			UserId:     event.Message.Sender.Uin,
			Message:    bot.MiraiMsgToProtoMsg(cli, event.Message.Elements),
			RawMessage: bot.MiraiMsgToRawMsg(cli, event.Message.Elements),
			Extra: map[string]string{
				"group_id": strconv.FormatInt(event.Message.GroupCode, 10),
			},
		},
	}
	bot.HandleEventFrame(cli, eventProto)
	return plugin.MessageIgnore
}

func ReportMemberPermissionChanged(cli *client.QQClient, event *client.MemberPermissionChangedEvent) int32 {
	eventProto := &onebot.Frame{
		FrameType: onebot.Frame_TGroupAdminNoticeEvent,
	}
	subType := "unset"
	if event.NewPermission == client.Administrator {
		subType = "set"
	}
	eventProto.Data = &onebot.Frame_GroupAdminNoticeEvent{
		GroupAdminNoticeEvent: &onebot.GroupAdminNoticeEvent{
			Time:       time.Now().Unix(),
			SelfId:     cli.Uin,
			PostType:   "notice",
			NoticeType: "group_admin",
			SubType:    subType,
			GroupId:    event.Group.Code,
			UserId:     event.Member.Uin,
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

func ReportGroupMute(cli *client.QQClient, event *client.GroupMuteEvent) int32 {
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

func ReportNewFriendRequest(cli *client.QQClient, event *client.NewFriendRequest) int32 {
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

func ReportUserJoinGroupRequest(cli *client.QQClient, event *client.UserJoinGroupRequest) int32 {
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

func ReportGroupInvitedRequest(cli *client.QQClient, event *client.GroupInvitedRequest) int32 {
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

func ReportGroupMessageRecalled(cli *client.QQClient, event *client.GroupMessageRecalledEvent) int32 {
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

func ReportFriendMessageRecalled(cli *client.QQClient, event *client.FriendMessageRecalledEvent) int32 {
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

func ReportNewFriendAdded(cli *client.QQClient, event *client.NewFriendEvent) int32 {
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
func ReportOfflineFile(cli *client.QQClient, event *client.OfflineFileEvent) int32 {
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
}
