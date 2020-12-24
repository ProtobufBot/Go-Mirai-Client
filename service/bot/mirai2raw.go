package bot

import (
	"fmt"
	"html"

	"github.com/Mrs4s/MiraiGo/message"
)

func MiraiMsgToRawMsg(messageChain []message.IMessageElement) string {
	result := ""
	for _, element := range messageChain {
		switch elem := element.(type) {
		case *message.TextElement:
			result += elem.Content
		case *message.AtElement:
			result += fmt.Sprintf(`<at qq="%d"/>`, elem.Target)
		case *message.ImageElement:
			result += fmt.Sprintf(`<image url="%s"/>`, html.EscapeString(elem.Url))
		case *message.FaceElement:
			result += fmt.Sprintf(`<face id="%d" name="%s"/>`, elem.Index, html.EscapeString(elem.Name))
		case *message.VoiceElement:
			result += fmt.Sprintf(`<voice url="%s"/>`, html.EscapeString(elem.Url))
		case *message.ServiceElement:
			result += fmt.Sprintf(`<service id="%d" content="%s" res_id="%s" sub_type="%s"/>`, elem.Id, html.EscapeString(elem.Content), elem.ResId, elem.SubType)
		case *message.LightAppElement:
			result += fmt.Sprintf(`<light_app content="%s"/>`, html.EscapeString(elem.Content))
		case *message.ShortVideoElement:
			result += fmt.Sprintf(`<video name="%s" url="%s"/>`, html.EscapeString(elem.Name), html.EscapeString(elem.Url))
		case *message.ReplyElement:
			result += fmt.Sprintf(`<reply time="%d" sender="%d" raw_message="%s" reply_seq="%d"/>`, elem.Time, elem.Sender, MiraiMsgToRawMsg(elem.Elements), elem.ReplySeq)
		}
	}
	return result
}
