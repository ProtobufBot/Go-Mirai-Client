# Go-Mirai-Client

用于收发QQ消息，并通过 websocket + protobuf 上报给 server 进行处理。

server端可以使用任意语言编写，通信协议：https://github.com/lz1998/onebot_idl

Java/Kotlin用户推荐使用 [spring-boot-starter](https://github.com/protobufbot/pbbot-spring-boot-starter)

支持发送的消息：文字、表情、图片、atQQ

## 验证码类型及处理方法

处理验证码时必须用到浏览器

### 短信验证码

右侧输入短信验证码提交

### 设备锁验证码

打开链接，扫码之后，右侧输入任意内容提交

在部分情况下，可以选择设备锁验证码(扫码)和短信验证码，默认选择扫码。如果需要默认选择短信，可以设置环境变量`SMS=1`

### 图形验证码

看图，在右侧输入验证码提交

### 滑块验证码

打开链接，F12打开浏览器开发者工具，选择Network，滑动滑块，出现cap_union_new_verify后，在response找ticket，复制到右侧输入框，提交

**重要：必须先打开 开发者工具-Network 之后滑动，否则可能看不到cap_union_new_verify**


## 自动登陆脚本

可以编写多个启动脚本自动启动不同的账号，不同账号PORT必须不同(或者写0表示随机)

如果需要登陆验证（图形/短信验证码等），必须使用浏览器访问`127.0.0.1:PORT`，PORT不能乱填

如果已经挂了很久，非常确定不会遇到验证码的时候，可以把环境变量的`PORT`设为0，表示使用随机端口。

### Windows

```shell
set UIN=机器人QQ号
set PASSWORD=机器人QQ密码
set PORT=9000
Go-Mirai-Client.exe
```

### Linux
```shell
chmod +x ./Go-Mirai-Client

export UIN=机器人QQ号
export PASSWORD=机器人QQ密码
export PORT=9000
./Go-Mirai-Client
```