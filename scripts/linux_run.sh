#!/bin/bash

read -p "Enter your bot qq:" uin
read -p "Enter your bot password:" pass

# 根据情况修改使用哪个
Go-Mirai-Client-linux-amd64 -uin "$uin" -pass "$pass"