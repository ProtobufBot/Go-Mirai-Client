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
			result += fmt.Sprintf("[@%d]", elem.Target)
		case *message.ImageElement:
			result += "[图片]"
		case *message.FaceElement:
			result += "[表情]"
		}
	}
	return result
}
