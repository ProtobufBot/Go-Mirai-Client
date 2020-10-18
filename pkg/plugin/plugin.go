package plugin

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

type (
	PrivateMessagePlugin        = func(*client.QQClient, *message.PrivateMessage) int32
	GroupMessagePlugin          = func(*client.QQClient, *message.GroupMessage) int32
	MemberJoinGroupPlugin       = func(*client.QQClient, *client.MemberJoinGroupEvent) int32
	MemberLeaveGroupPlugin      = func(*client.QQClient, *client.MemberLeaveGroupEvent) int32
	JoinGroupPlugin             = func(*client.QQClient, *client.GroupInfo) int32
	LeaveGroupPlugin            = func(*client.QQClient, *client.GroupLeaveEvent) int32
	NewFriendRequestPlugin      = func(*client.QQClient, *client.NewFriendRequest) int32
	UserJoinGroupRequestPlugin  = func(*client.QQClient, *client.UserJoinGroupRequest) int32
	GroupInvitedRequestPlugin   = func(*client.QQClient, *client.GroupInvitedRequest) int32
	GroupMessageRecalledPlugin  = func(*client.QQClient, *client.GroupMessageRecalledEvent) int32
	FriendMessageRecalledPlugin = func(*client.QQClient, *client.FriendMessageRecalledEvent) int32
	NewFriendAddedPlugin        = func(*client.QQClient, *client.NewFriendEvent) int32
)

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

func Serve(cli *client.QQClient) {
	cli.OnPrivateMessage(handlePrivateMessage)
	cli.OnGroupMessage(handleGroupMessage)
	cli.OnGroupMemberJoined(handleMemberJoinGroup)
	cli.OnGroupMemberLeaved(handleMemberLeaveGroup)
	cli.OnJoinGroup(handleJoinGroup)
	cli.OnLeaveGroup(handleLeaveGroup)
	cli.OnNewFriendRequest(handleNewFriendRequest)
	cli.OnUserWantJoinGroup(handleUserJoinGroupRequest)
	cli.OnGroupInvited(handleGroupInvitedRequest)
	cli.OnGroupMessageRecalled(handleGroupMessageRecalled)
	cli.OnFriendMessageRecalled(handleFriendMessageRecalled)
	cli.OnNewFriendAdded(handleNewFriendAdded)
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

func handlePrivateMessage(cli *client.QQClient, event *message.PrivateMessage) {
	SafeGo(func() {
		for _, plugin := range PrivateMessagePluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func handleGroupMessage(cli *client.QQClient, event *message.GroupMessage) {
	SafeGo(func() {
		for _, plugin := range GroupMessagePluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func handleMemberJoinGroup(cli *client.QQClient, event *client.MemberJoinGroupEvent) {
	SafeGo(func() {
		for _, plugin := range MemberJoinGroupPluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func handleMemberLeaveGroup(cli *client.QQClient, event *client.MemberLeaveGroupEvent) {
	SafeGo(func() {
		for _, plugin := range MemberLeaveGroupPluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func handleJoinGroup(cli *client.QQClient, event *client.GroupInfo) {
	SafeGo(func() {
		for _, plugin := range JoinGroupPluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func handleLeaveGroup(cli *client.QQClient, event *client.GroupLeaveEvent) {
	SafeGo(func() {
		for _, plugin := range LeaveGroupPluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func handleNewFriendRequest(cli *client.QQClient, event *client.NewFriendRequest) {
	SafeGo(func() {
		for _, plugin := range NewFriendRequestPluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func handleUserJoinGroupRequest(cli *client.QQClient, event *client.UserJoinGroupRequest) {
	SafeGo(func() {
		for _, plugin := range UserJoinGroupRequestPluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func handleGroupInvitedRequest(cli *client.QQClient, event *client.GroupInvitedRequest) {
	SafeGo(func() {
		for _, plugin := range GroupInvitedRequestPluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func handleGroupMessageRecalled(cli *client.QQClient, event *client.GroupMessageRecalledEvent) {
	SafeGo(func() {
		for _, plugin := range GroupMessageRecalledPluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func handleFriendMessageRecalled(cli *client.QQClient, event *client.FriendMessageRecalledEvent) {
	SafeGo(func() {
		for _, plugin := range FriendMessageRecalledPluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func handleNewFriendAdded(cli *client.QQClient, event *client.NewFriendEvent) {
	SafeGo(func() {
		for _, plugin := range NewFriendAddedPluginList {
			if result := plugin(cli, event); result == MessageBlock {
				break
			}
		}
	})
}

func SafeGo(fn func()) {
	go func() {
		defer func() {
			e := recover()
			if e != nil {
				fmt.Println(e)
			}
		}()
		fn()
	}()
}
