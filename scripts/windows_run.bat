@echo off

echo ================================================
echo               Go-Mirai-Client
echo You can modify this file to login automatically.
echo You can copy this file to manage many bots.
echo https://github.com/ProtobufBot/go-Mirai-Client
echo ================================================

set port=9000

set /p uin=QQ:
set /p pass=Password:
set /p port=Port(9000-60000):


Go-Mirai-Client-windows-amd64.exe -uin %uin% -pass %pass% -port %port%