package bot

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/ProtobufBot/Go-Mirai-Client/pkg/util"
	"github.com/ProtobufBot/Go-Mirai-Client/proto_gen/dto"
	"github.com/fanliao/go-promise"
	log "github.com/sirupsen/logrus"
)

var Captcha *dto.Captcha
var CaptchaPromise *promise.Promise

func ProcessLoginRsp(cli *client.QQClient, rsp *client.LoginResponse) (bool, error) {
	if rsp.Success {
		Captcha = nil
		return true, nil
	}
	switch rsp.Error {
	case client.SliderNeededError:
		if client.SystemDeviceInfo.Protocol == client.AndroidPhone {
			log.Warnf("警告: Android Phone 强制要求暂不支持的滑条验证码, 请开启设备锁或切换到Watch协议验证通过后再使用.")
			time.Sleep(5 * time.Second)
			os.Exit(0)
		}
		cli.AllowSlider = false
		cli.Disconnect()
		rsp, err := cli.Login()
		if err != nil {
			return false, err
		}
		return ProcessLoginRsp(cli, rsp)
	case client.NeedCaptcha:
		_ = ioutil.WriteFile("captcha.jpg", rsp.CaptchaImage, 0644)
		Captcha = &dto.Captcha{
			BotId:       cli.Uin,
			CaptchaType: dto.Captcha_PIC_CAPTCHA,
			Data:        &dto.Captcha_Image{Image: rsp.CaptchaImage},
		}
		CaptchaPromise = promise.NewPromise()
		result, err := CaptchaPromise.Get()
		text := result.(string)
		rsp, err := cli.SubmitCaptcha(strings.ReplaceAll(text, "\n", ""), rsp.CaptchaSign)
		util.DelFile("captcha.jpg")
		if err != nil {
			return false, fmt.Errorf("提交图形验证码错误")
		}
		return ProcessLoginRsp(cli, rsp)
	case client.SMSNeededError, client.SMSOrVerifyNeededError:
		if !cli.RequestSMS() {
			return false, fmt.Errorf("请求短信验证码错误，可能是太频繁")
		}
		Captcha = &dto.Captcha{
			BotId:       cli.Uin,
			CaptchaType: dto.Captcha_SMS,
			Data:        &dto.Captcha_Url{Url: rsp.SMSPhone},
		}
		CaptchaPromise = promise.NewPromise()
		result, err := CaptchaPromise.Get()
		if err != nil {
			return false, fmt.Errorf("提交短信验证码错误")
		}
		text := result.(string)
		rsp, err := cli.SubmitSMS(strings.ReplaceAll(strings.ReplaceAll(text, "\n", ""), "\r", ""))

		return ProcessLoginRsp(cli, rsp)
	case client.UnsafeDeviceError:
		Captcha = &dto.Captcha{
			BotId:       cli.Uin,
			CaptchaType: dto.Captcha_UNSAFE_DEVICE_LOGIN_VERIFY,
			Data:        &dto.Captcha_Url{Url: rsp.VerifyUrl},
		}
		CaptchaPromise = promise.NewPromise()
		_, err := CaptchaPromise.Get()
		rsp, err := cli.Login()
		if err != nil {
			return false, fmt.Errorf("设备锁验证/登陆错误")
		}
		return ProcessLoginRsp(cli, rsp)
	case client.OtherLoginError, client.UnknownLoginError:
		log.Errorf(rsp.ErrorMessage)
		return false, fmt.Errorf("遇到不可处理的登录错误")
	}
	return false, fmt.Errorf("process login error")
}
