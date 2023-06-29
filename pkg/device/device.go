package device

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ProtobufBot/Go-Mirai-Client/pkg/config"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/util"

	"github.com/Mrs4s/MiraiGo/binary"
	"github.com/Mrs4s/MiraiGo/client"
	log "github.com/sirupsen/logrus"
)

func RandDevice(randGen *rand.Rand) *client.DeviceInfo {
	device := &client.DeviceInfo{
		AndroidId:    []byte("MIRAI.123456.001"),
		Display:      []byte("GMC." + util.RandomStringRange(randGen, 6, "0123456789") + ".001"),
		Product:      []byte("gmc"),
		Device:       []byte("gmc"),
		Board:        []byte("gmc"),
		Brand:        []byte("pbbot"),
		Model:        []byte("gmc"),
		Bootloader:   []byte("unknown"),
		FingerPrint:  []byte("pbbot/gmc/gmc:10/PBBOT.200324.001/" + util.RandomStringRange(randGen, 7, "0123456789") + ":user/release-keys"),
		ProcVersion:  []byte("Linux 5.4.0-54-generic" + util.RandomString(randGen, 8) + " (android-build@gmail.com)"),
		BaseBand:     []byte{},
		SimInfo:      []byte("T-Mobile"),
		OSType:       []byte("android"),
		MacAddress:   []byte(fmt.Sprintf("EC:D0:9F:%s:%s:%s", util.RandomStringRange(randGen, 2, "02468ACD"), util.RandomStringRange(randGen, 2, "02468ACD"), util.RandomStringRange(randGen, 2, "02468ACD"))),
		IpAddress:    []byte{10, 0, 1, 3}, // 10.0.1.3
		WifiSSID:     []byte("TP-LINK-" + util.RandomStringRange(randGen, 6, "ABCDEF1234567890")),
		IMEI:         GenIMEI(randGen),
		APN:          []byte("wifi"),
		VendorName:   []byte("MIUI"),
		VendorOSName: []byte("gmc"),
		Protocol:     client.IPad,
		Version: &client.Version{
			Incremental: []byte("5891938"),
			Release:     []byte("10"),
			CodeName:    []byte("REL"),
			SDK:         29,
		},
	}
	device.WifiBSSID = device.MacAddress
	r := make([]byte, 16)
	randGen.Read(r)
	device.BootId = binary.GenUUID(r)
	randGen.Read(r)
	t := md5.Sum(r)
	device.IMSIMd5 = t[:]
	r = make([]byte, 8)
	randGen.Read(r)
	hex.Encode(device.AndroidId, r)
	GenNewGuid(device)
	GenNewTgtgtKey(randGen, device)
	return device
}

func GenIMEI(randGen *rand.Rand) string {
	sum := 0 // the control sum of digits
	var final strings.Builder

	for i := 0; i < 14; i++ { // generating all the base digits
		toAdd := randGen.Intn(10)
		if (i+1)%2 == 0 { // special proc for every 2nd one
			toAdd *= 2
			if toAdd >= 10 {
				toAdd = (toAdd % 10) + 1
			}
		}
		sum += toAdd
		final.WriteString(fmt.Sprintf("%d", toAdd)) // and even printing them here!
	}
	ctrlDigit := (sum * 9) % 10 // calculating the control digit
	final.WriteString(fmt.Sprintf("%d", ctrlDigit))
	return final.String()
}

func GenNewGuid(info *client.DeviceInfo) {
	t := md5.Sum(append(info.AndroidId, info.MacAddress...))
	info.Guid = t[:]
}

func GenNewTgtgtKey(randGen *rand.Rand, info *client.DeviceInfo) {
	r := make([]byte, 16)
	randGen.Read(r)
	h := md5.New()
	h.Write(r)
	h.Write(info.Guid)
	info.TgtgtKey = h.Sum(nil)
}

// GetDevice
// 如果设备文件夹不存在，自动创建文件夹
// 使用种子生成随机设备信息
// 如果已有设备文件，使用已有设备信息覆盖
// 存储设备信息到文件
func GetDevice(seed int64, clientProtocol int32) *client.DeviceInfo {
	randGen := rand.New(rand.NewSource(seed))
	if seed == 0 {
		randGen = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	// 默认 device/device-qq.json
	devicePath := path.Join("device", fmt.Sprintf("device-%d.json", seed))

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
	deviceInfo := RandDevice(randGen)

	deviceInfo.IpAddress = []byte{192, 168, 1, byte(100 + seed%100)}

	if util.PathExists(devicePath) {
		log.Infof("使用 %s 内的设备信息覆盖设备信息", devicePath)
		if err := deviceInfo.ReadJson([]byte(util.ReadAllText(devicePath))); err != nil {
			util.FatalError(fmt.Errorf("failed to load device info, err: %+v", err))
		}
	}

	if clientProtocol > 0 && clientProtocol < 7 {
		deviceInfo.Protocol = client.ClientProtocol(clientProtocol)
	} else {
		deviceInfo.Protocol = client.ClientProtocol(6)
	}

	GenNewGuid(deviceInfo)
	GenNewTgtgtKey(randGen, deviceInfo)

	log.Infof("保存设备信息到文件 %s", devicePath)
	err := ioutil.WriteFile(devicePath, deviceInfo.ToJson(), 0644)
	if err != nil {
		log.Warnf("写设备信息文件 %s 失败", devicePath)
	}
	return deviceInfo
}
