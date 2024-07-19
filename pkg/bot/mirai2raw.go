package bot

import (
	"fmt"
	"html"

	"github.com/2mf8/LagrangeGo/client"
	"github.com/2mf8/LagrangeGo/message"
)

func MiraiMsgToRawMsg(cli *client.QQClient, messageChain []message.IMessageElement) string {
	result := ""
	for _, element := range messageChain {
		switch elem := element.(type) {
		case *message.TextElement:
			result += elem.Content
		case *message.ImageElement:
			result += fmt.Sprintf(`<image image_id="%s" url="%s"/>`, html.EscapeString(elem.ImageId), html.EscapeString(elem.Url))
		case *message.FaceElement:
			result += fmt.Sprintf(`<face id="%d"/>`, elem.FaceID)
		case *message.VoiceElement:
			result += fmt.Sprintf(`<voice url="%s"/>`, html.EscapeString(elem.Url))
		case *message.ReplyElement:
			result += fmt.Sprintf(`<reply time="%d" sender="%d" raw_message="%s" reply_seq="%d"/>`, elem.Time, elem.SenderUin, html.EscapeString(MiraiMsgToRawMsg(cli, elem.Elements)), elem.ReplySeq)
		}
	}
	return result
}
