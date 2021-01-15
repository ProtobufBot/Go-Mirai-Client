package clz

import (
	"io"

	"github.com/Mrs4s/MiraiGo/message"
)

// 自定义类型

type VideoElement struct {
	message.ShortVideoElement
	UploadingCover io.ReadSeeker // 待上传的封面
	UploadingVideo io.ReadSeeker // 待上传的视频
}
