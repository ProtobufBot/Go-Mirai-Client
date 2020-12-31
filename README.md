# Go-Mirai-Client

[![QQ群](https://img.shields.io/static/v1?label=QQ%E7%BE%A4&message=335783090&color=blue)](https://jq.qq.com/?_wv=1027&k=B7Of3GMZ)

用于收发QQ消息，并通过 websocket + protobuf 上报给 server 进行处理。

server端可以使用任意语言编写，通信协议：https://github.com/lz1998/onebot_idl

Java/Kotlin用户推荐使用 [spring-boot-starter](https://github.com/protobufbot/pbbot-spring-boot-starter)

支持发送的消息：文字、表情、图片、atQQ

有问题发issue，或者进QQ群335783090

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

#### 方案A：自行抓包

由于滑块验证码和QQ本体的协议独立, 我们无法直接处理并提交. 需要在浏览器通过后抓包并获取 `Ticket` 提交.

该方案为具体的抓包教程, 如果您已经知道如何在浏览器中抓包. 可以略过接下来的文档并直接抓取 `cap_union_new_verify` 的返回值, 提取 `Ticket` 并在命令行提交.

首先打开滑块链接. 这里以 *Microsoft Edge* 浏览器为例, *Chrome* 同理. 

![image.png](https://i.loli.net/2020/12/27/otk9Hz7lBCaRFMV.png)

此时不要滑动验证码, 首先按下 `F12` (键盘右上角退格键上方) 打开 *开发者工具*

![image.png](https://i.loli.net/2020/12/27/JDioadLPwcKWpt1.png)

点击 `Network` 选项卡 (在某些浏览器它可能叫做 `网络`)

![image.png](https://i.loli.net/2020/12/27/qEzTB5jrDZUWSwp.png)

点开 `Filter` (箭头) 按钮以确定您能看到下面的工具栏, 勾选 `Preserve log`(红框)

此时可以滑动并通过验证码

![image.png](https://i.loli.net/2020/12/27/Id4hxzyDprQuF2G.png)

回到 *开发者工具*, 我们可以看到已经有了一个请求.

![image.png](https://i.loli.net/2020/12/27/3C6Y2XVKBRv1z9E.png)

此时如果有多个请求, 请不要慌张. 看到上面的 `Filter` 没? 此时在 `Filter` 输入框中输入 `cap_union_new`, 就应该只剩一个请求了.

然后点击该请求. 点开 `Preview` 选项卡 (箭头):  

![image.png](https://i.loli.net/2020/12/27/P1VtxRWpjY8524Z.png)

此时就能看到一个标准的 `JSON`, 复制 `ticket` 字段并回到 验证码处理 粘贴. 即可通过滑块验证.

#### 方案B

下载专用软件进行验证

TODO


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
