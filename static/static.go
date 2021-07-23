package static

import (
	"embed"
	"io/fs"
)

// 需要把前端文件放在static文件夹

//go:embed static
var staticFs embed.FS

func MustGetStatic() fs.FS {
	f, err := fs.Sub(staticFs, "static")
	if err != nil {
		panic(err)
	}
	return f
}
