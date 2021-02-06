package bot

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/ProtobufBot/Go-Mirai-Client/config"
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
	if rsp.Error == client.SMSOrVerifyNeededError {
		if config.Conf.SMS {
			rsp.Error = client.SMSNeededError
		} else {
			rsp.Error = client.UnsafeDeviceError
		}
	}
	log.Infof("验证码处理页面: http://localhost:%s/", config.Conf.Port)
	switch rsp.Error {
	case client.SliderNeededError:
		log.Infof("遇到滑块验证码，根据README提示操作 https://github.com/protobufbot/Go-Mirai-Client (顺便star)")
		Captcha = &dto.Captcha{
			BotId:       cli.Uin,
			CaptchaType: dto.Captcha_SLIDER_CAPTCHA,
			Data:        &dto.Captcha_Url{Url: rsp.VerifyUrl},
		}
		CaptchaPromise = promise.NewPromise()
		result, err := CaptchaPromise.Get()
		if err != nil {
			return false, fmt.Errorf("提交ticket错误")
		}
		text := result.(string)
		rsp, err := cli.SubmitTicket(text)
		if err != nil {
			return false, err
		}
		return ProcessLoginRsp(cli, rsp)
	case client.NeedCaptcha:
		log.Infof("遇到图形验证码，根据README提示操作 https://github.com/protobufbot/Go-Mirai-Client (顺便star)")
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
	case client.SMSNeededError:
		log.Infof("遇到短信验证码，根据README提示操作 https://github.com/protobufbot/Go-Mirai-Client (顺便star)")
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
		log.Infof("遇到设备锁扫码验证码，根据README提示操作 https://github.com/protobufbot/Go-Mirai-Client (顺便star)")
		log.Info("设置环境变量 SMS = 1 可优先使用短信验证码")
		Captcha = &dto.Captcha{
			BotId:       cli.Uin,
			CaptchaType: dto.Captcha_UNSAFE_DEVICE_LOGIN_VERIFY,
			Data:        &dto.Captcha_Url{Url: rsp.VerifyUrl},
		}
		CaptchaPromise = promise.NewPromise()
		_, err := CaptchaPromise.Get()
		cli.Disconnect()
		time.Sleep(3 * time.Second)
		rsp, err := cli.Login()
		if err != nil {
			return false, fmt.Errorf("设备锁验证/登陆错误")
		}
		return ProcessLoginRsp(cli, rsp)
	case client.OtherLoginError, client.UnknownLoginError:
		//log.Errorf(rsp.ErrorMessage)
		log.Warnf("登陆失败，建议开启/关闭设备锁后重试，或删除device-<QQ>.json文件后再次尝试")
		msg := rsp.ErrorMessage
		if strings.Contains(msg, "版本") {
			log.Errorf("密码错误或账号被冻结")
		}
		if strings.Contains(msg, "上网环境") {
			log.Errorf("当前上网环境异常. 更换服务器并重试")
		}
		return false, fmt.Errorf("遇到不可处理的登录错误")
	}
	return false, fmt.Errorf("process login error")
}
