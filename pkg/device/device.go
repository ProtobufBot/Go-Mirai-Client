package device

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/2mf8/Go-Lagrange-Client/pkg/config"
	"github.com/2mf8/Go-Lagrange-Client/pkg/util"
	"github.com/2mf8/LagrangeGo/client/auth"
	log "github.com/sirupsen/logrus"
)

// GetDevice
// 如果设备文件夹不存在，自动创建文件夹
// 使用种子生成随机设备信息
// 如果已有设备文件，使用已有设备信息覆盖
// 存储设备信息到文件
func GetDevice(seed int64) *auth.DeviceInfo {
	// 默认 device/device-qq.json
	devicePath := path.Join("device", fmt.Sprintf("%d.json", seed))

	// 优先使用参数目录
	if config.Device != "" {
		devicePath = config.Device
	}

	deviceDir := path.Dir(devicePath)
	if !util.PathExists(deviceDir) {
		log.Infof("%+v 目录不存在，自动创建", deviceDir)
		if err := os.MkdirAll(deviceDir, 0777); err != nil {
			log.Warnf("failed to mkdir deviceDir, err: %+v", err)
		}
	}

	log.Info("生成随机设备信息")
	deviceInfo := auth.NewDeviceInfo(int(seed))

	if util.PathExists(devicePath) {
		log.Infof("使用 %s 内的设备信息覆盖设备信息", devicePath)
		fi, err := os.ReadFile(devicePath)
		if err != nil {
			util.FatalError(fmt.Errorf("failed to read device info, err: %+v", err))
		}
		err = json.Unmarshal(fi, deviceInfo)
		if err != nil {
			util.FatalError(fmt.Errorf("failed to load device info, err: %+v", err))
		}
	}

	log.Infof("保存设备信息到文件 %s", devicePath)
	data, err := json.Marshal(deviceInfo)
	if err != nil {
		log.Warnf("JSON 化设备信息文件 %s 失败", devicePath)
	}
	err = os.WriteFile(devicePath, data, 0644)
	if err != nil {
		log.Warnf("写设备信息文件 %s 失败", devicePath)
	}
	return deviceInfo
}