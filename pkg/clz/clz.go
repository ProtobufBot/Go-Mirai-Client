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

type LocalImageElement struct {
	message.ImageElement
	Stream   io.ReadSeeker
	Tp       string // 类型 flash/show
	EffectId int32 // show的特效id，范围40000-40005
}

type GiftElement struct {
	Target int64
	GiftId message.GroupGift
}

func (g *GiftElement) Type() message.ElementType {
	return message.At
}
