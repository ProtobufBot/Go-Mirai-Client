# Go-Mirai-Client

用于收发QQ消息，并通过 websocket + protobuf 上报给 server 进行处理。

server端可以使用任意语言编写，通信协议：https://github.com/lz1998/onebot_idl

Java/Kotlin用户推荐使用 [spring-boot-starter](https://github.com/protobufbot/pbbot-spring-boot-starter)

支持发送的消息：文字、表情、图片、atQQ


## 自动登陆脚本

可以编写多个启动脚本自动启动不同的账号，不同账号PORT必须不同(或者写0表示随机)

如果需要登陆验证（图形/短信验证码等），必须使用浏览器访问`127.0.0.1:PORT`，PORT不能乱填

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