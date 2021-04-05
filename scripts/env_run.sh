#!/bin/bash

cmd="gmc "

if [ $UIN ];then
  cmd+="-uin $UIN "
fi

if [ $PASS ];then
  cmd+="-pass $PASS "
fi

if [ $PORT ];then
  cmd+="-port $PORT "
fi

if [ $WS_URL ];then
  cmd+="-ws_url $WS_URL "
fi

if [ $SMS ];then
  cmd+="-sms $SMS "
fi

if [ $DEVICE ];then
  cmd+="-device $DEVICE "
fi

eval $cmd