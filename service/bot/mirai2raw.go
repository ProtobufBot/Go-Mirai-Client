package bot

import (
	"fmt"
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
			result += fmt.Sprintf(`<image file="%s" url="%s"/>`, elem.Url, elem.Url)
		case *message.FaceElement:
			result += fmt.Sprintf(`<face id="%d" name="%s"/>`, elem.Index, elem.Name)
		case *message.VoiceElement:
			result += fmt.Sprintf(`<voice file="%s" url="%s"/>`, elem.Url, elem.Url)
		}
	}
	return result
}
