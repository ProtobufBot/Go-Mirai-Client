# Go-Mirai-Client

[![QQ群](https://img.shields.io/static/v1?label=QQ%E7%BE%A4&message=335783090&color=blue)](https://jq.qq.com/?_wv=1027&k=B7Of3GMZ)

用于收发QQ消息，并通过 websocket + protobuf 上报给 server 进行处理。

server端可以使用任意语言编写，通信协议：https://github.com/lz1998/onebot_idl

Java/Kotlin用户推荐使用 [spring-boot-starter](https://github.com/protobufbot/pbbot-spring-boot-starter)

支持发送的消息：文字、表情、图片、atQQ

有问题发issue，或者进QQ群335783090

## 自动登陆脚本

可以编写多个启动脚本自动启动不同的账号，不同账号PORT必须不同(或者写0表示随机)

如果需要登陆验证（图形/短信验证码等），必须使用浏览器访问`127.0.0.1:PORT`，PORT不能乱填

如果已经挂了很久，非常确定不会遇到验证码的时候，可以把参数的`port`设为0，表示使用随机端口。

### 参数

```shell
Usage of GMC:
  -uin int
        机器人QQ
  -pass string
        机器人密码
  -port int
        http管理端口(默认 9000), 0表示随机, 如果不需要处理验证码, 可以随便填
  -sms bool
        登录优先使用短信验证
  -ws_url string
        消息处理websocket服务器地址
  -device string
        设备文件位置
  -help
        帮助
```

### Windows

创建一个文件，后缀为`.bat`，写入以下内容，双击运行

```shell
Go-Mirai-Client.exe -uin <机器人QQ> -pass <机器人密码> -port <HTTP端口> -device <设备信息位置> -ws_url <消息处理器地址> -sms <是否优先短信登录>
```

### Linux

创建一个文件，后缀为`.sh`，写入以下内容，添加执行权限，运行

```shell
chmod +x ./Go-Mirai-Client

./Go-Mirai-Client -uin <机器人QQ> -pass <机器人密码> -port <HTTP端口> -device <设备信息位置> -ws_url <消息处理器地址> -sms <是否优先短信登录>
```

### Docker

```shell
docker run -it \
--name=gmc \
-p 9000:9000 \
-e UIN=<账号> \
-e PASS=<密码> \
-e WS_URL=<WebSocket地址> \
-e DEVICE=/deivce/123.json \
-v <设备文件目录>:/deivce \
lz1998/gmc:0.1.11
```

## 修改协议

修改device文件的protocol，对应关系如下：

- IPad: 0
- AndroidPhone: 1
- AndroidWatch: 2
- MacOS: 3
- 企点: 4

输入其他数字默认表示IPad

## 验证码类型及处理方法

处理验证码时必须用到浏览器

### 短信验证码

右侧输入短信验证码提交

### 设备锁验证码

复制链接到手机打开处理。

在部分情况下，可以选择设备锁验证码(扫码)和短信验证码，默认选择扫码。如果需要默认选择短信，可以使用参数`-sms`

### 图形验证码

看图，在右侧输入验证码提交

### 滑块验证码

使用[滑块验证助手](https://github.com/mzdluo123/TxCaptchaHelper/releases)