package cache

import (
	"github.com/golang/groupcache/lru"
)

var PrivateMessageLru = lru.New(512)
var GroupMessageLru = lru.New(2048)
var FriendRequestLru = lru.New(128)
var GroupRequestLru = lru.New(128)
