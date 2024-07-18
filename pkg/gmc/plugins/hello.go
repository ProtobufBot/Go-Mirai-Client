package plugins

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/2mf8/Go-Lagrange-Client/pkg/bot"
	"github.com/2mf8/Go-Lagrange-Client/pkg/plugin"
	"github.com/2mf8/LagrangeGo/client"
	"github.com/2mf8/LagrangeGo/message"
	log "github.com/sirupsen/logrus"
)

func HelloPrivateMessage(cli *client.QQClient, event *message.PrivateMessage) int32 {
	if bot.MiraiMsgToRawMsg(cli, event.Elements) != "hi" {
		return plugin.MessageIgnore
	}
	elem := &message.SendingMessage{
		Elements: []message.IMessageElement{
			&message.TextElement{Content: "hello"},
		},
	}
	cli.SendPrivateMessage(event.Sender.Uin, elem.Elements)
	return plugin.MessageIgnore
}

func HelloGroupMessage(cli *client.QQClient, event *message.GroupMessage) int32 {
	if bot.MiraiMsgToRawMsg(cli, event.Elements) != "hi" {
		return plugin.MessageIgnore
	}
	resp, err := http.Get("https://www.2mf8.cn/static/image/cube3/b1.png")
	defer resp.Body.Close()
	fmt.Println(err)
	if err != nil {
		return plugin.MessageIgnore
	}
	imo, err := io.ReadAll(resp.Body)
	fmt.Println(err)
	if err != nil {
		return plugin.MessageIgnore
	}
	filename := fmt.Sprintf("%v.png", time.Now().UnixMicro())
	err = os.WriteFile(filename, imo, 0666)
	fmt.Println(err)
	if err != nil {
		return plugin.MessageIgnore
	}
	f, err := os.Open(filename)
	fmt.Println(err)
	if err != nil {
		return plugin.MessageIgnore
	}
	ir, err := cli.ImageUploadGroup(event.GroupUin, message.NewStreamImage(f))
	fmt.Println(err)
	if err != nil {
		return plugin.MessageIgnore
	}
	elem := &message.SendingMessage{}
	elem.Elements = append(elem.Elements, ir)
	elem.Elements = append(elem.Elements, &message.TextElement{
		Content: "测试成功",
	})
	r, e := cli.SendGroupMessage(event.GroupUin, elem.Elements)
	log.Warn(r, e)
	return plugin.MessageIgnore
}
