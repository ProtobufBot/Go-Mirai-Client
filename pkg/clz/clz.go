package clz

import (
	"io"

	"github.com/Mrs4s/MiraiGo/message"
)

// 自定义类型

type MyVideoElement struct {
	message.ShortVideoElement
	CoverUrl       string        // 仅用于发送时日志展示
	UploadingCover io.ReadSeeker // 待上传的封面 发送时需要
	UploadingVideo io.ReadSeeker // 待上传的视频 发送时需要
}
