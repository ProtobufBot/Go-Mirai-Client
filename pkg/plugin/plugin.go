package plugin

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

type (
	PrivateMessagePlugin   = func(*client.QQClient, *message.PrivateMessage) int32
	GroupMessagePlugin     = func(*client.QQClient, *message.GroupMessage) int32
	MemberJoinGroupPlugin  = func(*client.QQClient, *client.MemberJoinGroupEvent) int32
	MemberLeaveGroupPlugin = func(*client.QQClient, *client.MemberLeaveGroupEvent) int32
	JoinGroupPlugin        = func(*client.QQClient, *client.GroupInfo) int32
	LeaveGroupPlugin       = func(*client.QQClient, *client.GroupLeaveEvent) int32
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

func Serve(cli *client.QQClient) {
	cli.OnPrivateMessage(handlePrivateMessage)
	cli.OnGroupMessage(handleGroupMessage)
	cli.OnGroupMemberJoined(handleMemberJoinGroup)
	cli.OnGroupMemberLeaved(handleMemberLeaveGroup)
	cli.OnJoinGroup(handleJoinGroup)
	cli.OnLeaveGroup(handleLeaveGroup)
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
