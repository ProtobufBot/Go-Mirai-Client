# Go-Lagrange-Client

# Go-Mirai-Client 已经停更，新项目已迁移至 [Go-Lagrange-Client](https://github.com/2mf8/Go-Lagrange-Client)

## 默认端口9000，若启动失败，将依次在9001-9020端口中选择可用的端口启动。

用于收发QQ消息，并通过 websocket + protobuf 或 websocket + json 上报给 server 进行处理。

已支持与 `OneBot V11` 协议的服务端通信, 使用前需要选用 OneBot V11 协议

Golang 推荐使用 [GoneBot](https://github.com/2mf8/GoneBot)
TypeScript / JavaScript 推荐使用 [ToneBot](https://github.com/2mf8/ToneBot)

可以使用任意语言编写websocket server实现通信，协议：[onebot_glc](https://github.com/2mf8/onebot_glc)

有问题发issue，或者进QQ群 `901125207`

支持的开发语言(需要根据协议修改)：[Java/Kotlin](https://github.com/protobufbot/spring-mirai-server) , [JavaScript](https://github.com/2mf8/TSPbBot) , [TypeScript](https://github.com/2mf8/TSPbBot/blob/master/src/demo/index.ts) , [Python](https://github.com/PHIKN1GHT/pypbbot/tree/main/pypbbot_examples) , [Golang](https://github.com/2mf8/GoPbBot/blob/master/test/bot_test.go) , [C/C++](https://github.com/ProtobufBot/cpp-pbbot/blob/main/src/event_handler/event_handler.cpp) , [易语言](https://github.com/protobufbot/pbbot_e_sdk) 。详情查看 [Protobufbot](https://github.com/ProtobufBot/ProtobufBot) 。

## 使用说明

1. 启动程序
    - 用户在 [Releases](https://github.com/ProtobufBot/Go-Mirai-Client/releases) 下载适合自己的版本运行，然后手动打开浏览器地址`http://localhost:9000/`，Linux服务器可以远程访问`http://<服务器地址>:9000`。

2. 创建机器人
    - 建议选择扫码创建，使用**机器人账号**直接扫码，点击确认后登录。
    - 每次登录**必须**使用相同随机种子（数字），方便后续 `session` 登录。（建议使用账号作为随机种子）

3. 配置消息处理器
    - 在首次启动自动生成的`default.json`中配置服务器URL，修改后重启生效。
    - 如果使用其他人编写的程序，建议把`default.json`打包在一起发送给用户。

## 多插件支持

支持多插件，且每个插件URL可以配置多个作为候选项

cd service/glc && go run *.go
