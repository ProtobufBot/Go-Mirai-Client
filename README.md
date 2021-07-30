# Go-Mirai-Client

[![QQ群](https://img.shields.io/static/v1?label=QQ%E7%BE%A4&message=335783090&color=blue)](https://jq.qq.com/?_wv=1027&k=B7Of3GMZ)

用于收发QQ消息，并通过 websocket + protobuf 上报给 server 进行处理。

支持的开发语言：[Java/Kotlin](https://github.com/protobufbot/spring-mirai-server) , [JavaScript](https://github.com/ProtobufBot/js-pbbot/blob/master/example/src/index.js) , [Python](https://github.com/PHIKN1GHT/pypbbot/tree/main/pypbbot_examples) , [Golang](https://github.com/ProtobufBot/go-pbbot/blob/master/test/bot_test.go) , [C/C++](https://github.com/ProtobufBot/cpp-pbbot/blob/main/src/event_handler/event_handler.cpp) , [易语言](https://github.com/protobufbot/pbbot_e_sdk) 。详情查看 [Protobufbot](https://github.com/ProtobufBot/ProtobufBot) 。

可以使用其他任意语言编写websocket server实现通信，协议：[onebot_idl](https://github.com/lz1998/onebot_idl)

有问题发issue，或者进QQ群335783090

## 使用说明

1. 启动程序
    - Windows 非专业用户在 [Releases](https://github.com/ProtobufBot/Go-Mirai-Client/releases) 下载带有`lorca`
      的版本，启动时会自动打开UI界面（需要Edge/Chrome浏览器，安装在默认位置）。
    - 专业用户可以下载不带有`lorca`的版本，手动打开浏览器地址`http://localhost:9000/`，端口号可以通过`-port 9000`
      参数修改，Linux服务器可以远程访问`http://<服务器地址>:9000`。

2. 创建机器人
    - 建议选择扫码创建，使用**机器人账号**直接扫码，点击确认后登录。
    - 使用密码创建可能处理验证码。
    - 每次登录**必须**使用相同随机种子（数字），否则容易冻结。（建议使用账号作为随机种子）

3. 配置消息处理器
    - 在首次启动自动生成的`gmc_config.json`中配置服务器URL，修改后重启生效。
    - 如果使用其他人编写的程序，建议把`gmc_config.json`打包在一起发送给用户。

## 验证码类型及处理方法

使用密码登录会遇到验证码，点击机器人下方图标处理验证码，处理验证码时必须用到浏览器。

1. 设备锁验证码：复制链接到手机打开处理，可能需要扫码，如果添加参数`-sms`会优先使用短信验证码。
2. 短信验证码：直接输入短信内容提交。
3. 滑块验证码：使用[滑块验证助手](https://github.com/mzdluo123/TxCaptchaHelper/releases)。

## 运行参数

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

## 自动登陆

有2种方式可以实现自动登录。

1. 发送HTTP请求自动登录 支持多账号
2. 使用运行参数自动登录 只能单账号

### 发送HTTP请求自动登录

启动程序后，通过编写脚本，发送请求实现自动登录`POST http://localhost/bot/create/v1/`

```json
{
  "bot_id": 123,
  "password": "xxx",
  "device_seed": 123
}
```

### 使用运行参数自动登录

#### Windows

创建一个文件，后缀为`.bat`，写入以下内容，双击运行

```shell
Go-Mirai-Client.exe -uin <机器人QQ> -pass <机器人密码> -port <HTTP端口> -device <设备信息位置> -ws_url <消息处理器地址> -sms <是否优先短信登录>
```

#### Linux

创建一个文件，后缀为`.sh`，写入以下内容，添加执行权限，运行

```shell
chmod +x ./Go-Mirai-Client

./Go-Mirai-Client -uin <机器人QQ> -pass <机器人密码> -port <HTTP端口> -device <设备信息位置> -ws_url <消息处理器地址> -sms <是否优先短信登录>
```

## 多开

每次启动必须使用不同端口，默认使用9000端口。可以通过指定参数`-port 9000`修改端口，端口设置为0表示随机端口。

## Docker

```shell
docker run -it \
--name=gmc \
-p 9000:9000 \
-e UIN=<账号> \
-e PASS=<密码> \
-e WS_URL=<WebSocket地址> \
-e DEVICE=/deivce/123.json \
-v <设备文件目录>:/deivce \
lz1998/gmc:0.1.19
```

## 修改协议

修改device文件的protocol，对应关系如下：

- IPad: 0
- AndroidPhone: 1
- AndroidWatch: 2
- MacOS: 3
- 企点: 4

输入其他数字默认表示IPad
