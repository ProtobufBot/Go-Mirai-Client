#!/bin/sh

CMD="gmc"

if [ $UIN ];then
  CMD="$CMD -uin $UIN"
fi

if [ $PASS ];then
  CMD="$CMD -pass $PASS"
fi

if [ $PORT ];then
  CMD="$CMD -port $PORT"
fi

if [ $WS_URL ];then
  CMD="$CMD -ws_url $WS_URL"
fi

if [ $SMS ];then
  CMD="$CMD -sms $SMS"
fi

if [ $DEVICE ];then
  CMD="$CMD -device $DEVICE"
fi

if [ $AUTH ];then
  CMD="$CMD -auth $AUTH"
fi
echo $CMD
eval $CMD