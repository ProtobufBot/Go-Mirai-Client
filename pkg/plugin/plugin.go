package plugin

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

type (
	PrivateMessagePlugin = func(*client.QQClient, *message.PrivateMessage) int32
	GroupMessagePlugin   = func(*client.QQClient, *message.GroupMessage) int32
)

const (
	MessageIgnore = 0
	MessageBlock  = 1
)

var PrivateMessagePluginList = make([]PrivateMessagePlugin, 0)
var GroupMessagePluginList = make([]GroupMessagePlugin, 0)

func Serve(cli *client.QQClient) {
	cli.OnPrivateMessage(handlePrivateMessage)
	cli.OnGroupMessage(handleGroupMessage)

}

// 添加私聊消息插件
func AddPrivateMessagePlugin(plugin PrivateMessagePlugin) {
	PrivateMessagePluginList = append(PrivateMessagePluginList, plugin)
}

// 添加群聊消息插件
func AddGroupMessagePlugin(plugin GroupMessagePlugin) {
	GroupMessagePluginList = append(GroupMessagePluginList, plugin)
}

func handlePrivateMessage(cli *client.QQClient, msg *message.PrivateMessage) {
	SafeGo(func() {
		for _, plugin := range PrivateMessagePluginList {
			if result := plugin(cli, msg); result == MessageBlock {
				break
			}
		}
	})
}

func handleGroupMessage(cli *client.QQClient, msg *message.GroupMessage) {
	SafeGo(func() {
		for _, plugin := range GroupMessagePluginList {
			if result := plugin(cli, msg); result == MessageBlock {
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
