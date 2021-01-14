package clz

import "github.com/Mrs4s/MiraiGo/message"

// 自定义类型

type VideoElement struct {
	message.ShortVideoElement
	UploadingCoverBytes []byte // 待上传的封面
	UploadingVideoBytes []byte // 待上传的视频
}
