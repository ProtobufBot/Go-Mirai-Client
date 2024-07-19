package plugin

import (
	"github.com/2mf8/LagrangeGo/client"
	"github.com/2mf8/LagrangeGo/client/event"
	"github.com/2mf8/LagrangeGo/message"
	"github.com/2mf8/Go-Lagrange-Client/pkg/util"
)

type (
	PrivateMessagePlugin        = func(*client.QQClient, *message.PrivateMessage) int32
	GroupMessagePlugin          = func(*client.QQClient, *message.GroupMessage) int32
	MemberJoinGroupPlugin       = func(*client.QQClient, *event.GroupMemberIncrease) int32
	MemberLeaveGroupPlugin      = func(*client.QQClient, *event.GroupMemberDecrease) int32
	JoinGroupPlugin             = func(*client.QQClient, *event.GroupMemberJoinRequest) int32
	LeaveGroupPlugin            = func(*client.QQClient, *event.GroupMemberDecrease) int32
	NewFriendRequestPlugin      = func(*client.QQClient, *event.NewFriendRequest) int32
	UserJoinGroupRequestPlugin  = func(*client.QQClient, *event.GroupMemberIncrease) int32
	GroupInvitedRequestPlugin   = func(*client.QQClient, *event.GroupInvite) int32
	GroupMessageRecalledPlugin  = func(*client.QQClient, *event.GroupRecall) int32
	FriendMessageRecalledPlugin = func(*client.QQClient, *event.FriendRecall) int32
	NewFriendAddedPlugin        = func(*client.QQClient, *event.NewFriendRequest) int32
	GroupMutePlugin             = func(*client.QQClient, *event.GroupMute) int32
)

var eclient *client.QQClient

const (
	MessageIgnore = 0
	MessageBlock  = 1
)

var PrivateMessagePluginList = make([]PrivateMessagePlugin, 0)
var GroupMessagePluginList = make([]GroupMessagePlugin, 0)
var MemberJoinGroupPluginList = make([]MemberJoinGroupPlugin, 0)
var MemberLeaveGroupPluginList = make([]MemberLeaveGroupPlugin, 0)
var JoinGroupPluginList = make([]JoinGroupPlugin, 0)
var LeaveGroupPluginList = make([]LeaveGroupPlugin, 0)
var NewFriendRequestPluginList = make([]NewFriendRequestPlugin, 0)
var UserJoinGroupRequestPluginList = make([]UserJoinGroupRequestPlugin, 0)
var GroupInvitedRequestPluginList = make([]GroupInvitedRequestPlugin, 0)
var GroupMessageRecalledPluginList = make([]GroupMessageRecalledPlugin, 0)
var FriendMessageRecalledPluginList = make([]FriendMessageRecalledPlugin, 0)
var NewFriendAddedPluginList = make([]NewFriendAddedPlugin, 0)
var GroupMutePluginList = make([]GroupMutePlugin, 0)

func Serve(cli *client.QQClient) {
	cli.PrivateMessageEvent.Subscribe(handlePrivateMessage)
	cli.GroupMessageEvent.Subscribe(handleGroupMessage)
	cli.GroupMemberJoinEvent.Subscribe(handleMemberJoinGroup)
	cli.GroupMemberLeaveEvent.Subscribe(handleMemberLeaveGroup)
	cli.GroupMemberLeaveEvent.Subscribe(handleLeaveGroup)
	cli.NewFriendRequestEvent.Subscribe(handleNewFriendRequest)
	cli.GroupInvitedEvent.Subscribe(handleGroupInvitedRequest)
	cli.GroupRecallEvent.Subscribe(handleGroupMessageRecalled)
	cli.FriendRecallEvent.Subscribe(handleFriendMessageRecalled)
	cli.GroupMuteEvent.Subscribe(handleGroupMute)
}

// 添加私聊消息插件
func AddPrivateMessagePlugin(plugin PrivateMessagePlugin) {
	PrivateMessagePluginList = append(PrivateMessagePluginList, plugin)
}

// 添加群聊消息插件
func AddGroupMessagePlugin(plugin GroupMessagePlugin) {
	GroupMessagePluginList = append(GroupMessagePluginList, plugin)
}

// 添加群成员加入插件
func AddMemberJoinGroupPlugin(plugin MemberJoinGroupPlugin) {
	MemberJoinGroupPluginList = append(MemberJoinGroupPluginList, plugin)
}

// 添加群成员离开插件
func AddMemberLeaveGroupPlugin(plugin MemberLeaveGroupPlugin) {
	MemberLeaveGroupPluginList = append(MemberLeaveGroupPluginList, plugin)
}

// 添加机器人进群插件
func AddJoinGroupPlugin(plugin JoinGroupPlugin) {
	JoinGroupPluginList = append(JoinGroupPluginList, plugin)
}

// 添加机器人离开群插件
func AddLeaveGroupPlugin(plugin LeaveGroupPlugin) {
	LeaveGroupPluginList = append(LeaveGroupPluginList, plugin)
}

// 添加好友请求处理插件
func AddNewFriendRequestPlugin(plugin NewFriendRequestPlugin) {
	NewFriendRequestPluginList = append(NewFriendRequestPluginList, plugin)
}

// 添加加群请求处理插件
func AddUserJoinGroupRequestPlugin(plugin UserJoinGroupRequestPlugin) {
	UserJoinGroupRequestPluginList = append(UserJoinGroupRequestPluginList, plugin)
}

// 添加机器人被邀请处理插件
func AddGroupInvitedRequestPlugin(plugin GroupInvitedRequestPlugin) {
	GroupInvitedRequestPluginList = append(GroupInvitedRequestPluginList, plugin)
}

// 添加群消息撤回处理插件
func AddGroupMessageRecalledPlugin(plugin GroupMessageRecalledPlugin) {
	GroupMessageRecalledPluginList = append(GroupMessageRecalledPluginList, plugin)
}

// 添加好友消息撤回处理插件
func AddFriendMessageRecalledPlugin(plugin FriendMessageRecalledPlugin) {
	FriendMessageRecalledPluginList = append(FriendMessageRecalledPluginList, plugin)
}

// 添加好友添加处理插件
func AddNewFriendAddedPlugin(plugin NewFriendAddedPlugin) {
	NewFriendAddedPluginList = append(NewFriendAddedPluginList, plugin)
}

// 添加群成员被禁言插件
func AddGroupMutePlugin(plugin GroupMutePlugin) {
	GroupMutePluginList = append(GroupMutePluginList, plugin)
}

func handlePrivateMessage(cli *client.QQClient, event *message.PrivateMessage) {
	util.SafeGo(func() {
		for _, plugin := range PrivateMessagePluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func handleGroupMessage(cli *client.QQClient, event *message.GroupMessage) {
	util.SafeGo(func() {
		for _, plugin := range GroupMessagePluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func handleMemberJoinGroup(cli *client.QQClient, event *event.GroupMemberIncrease) {
	util.SafeGo(func() {
		for _, plugin := range MemberJoinGroupPluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func handleMemberLeaveGroup(cli *client.QQClient, event *event.GroupMemberDecrease) {
	util.SafeGo(func() {
		for _, plugin := range MemberLeaveGroupPluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func handleLeaveGroup(cli *client.QQClient, event *event.GroupMemberDecrease) {
	util.SafeGo(func() {
		for _, plugin := range LeaveGroupPluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func handleNewFriendRequest(cli *client.QQClient, event *event.NewFriendRequest) {
	util.SafeGo(func() {
		for _, plugin := range NewFriendRequestPluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func handleGroupInvitedRequest(cli *client.QQClient, event *event.GroupInvite) {
	util.SafeGo(func() {
		for _, plugin := range GroupInvitedRequestPluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func handleGroupMessageRecalled(cli *client.QQClient, event *event.GroupRecall) {
	util.SafeGo(func() {
		for _, plugin := range GroupMessageRecalledPluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func handleFriendMessageRecalled(cli *client.QQClient, event *event.FriendRecall) {
	util.SafeGo(func() {
		for _, plugin := range FriendMessageRecalledPluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func handleGroupMute(cli *client.QQClient, event *event.GroupMute) {
	util.SafeGo(func() {
		for _, plugin := range GroupMutePluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}
