# Go-Mirai-Client

用于收发QQ消息，并通过 websocket + protobuf 上报给 server 进行处理。

server端可以使用任意语言编写，通信协议：https://github.com/lz1998/onebot_idl

Java/Kotlin用户推荐使用 [spring-boot-starter](https://github.com/protobufbot/pbbot-spring-boot-starter)

支持发送的消息：文字、表情、图片、atQQ
