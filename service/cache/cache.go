package cache

import (
	"sync/atomic"

	"github.com/golang/groupcache/lru"
)

var GlobalSeq int32 = 0

func NextGlobalSeq() int32 {
	return atomic.AddInt32(&GlobalSeq, 1)
}

var PrivateMessageLru = lru.New(512)
var GroupMessageLru = lru.New(2048)
var FriendRequestLru = lru.New(128)
var GroupRequestLru = lru.New(128)
