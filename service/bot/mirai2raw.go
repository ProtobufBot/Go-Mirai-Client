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
		}
	}
	return result
}
