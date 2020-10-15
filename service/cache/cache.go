package cache

import (
	"github.com/golang/groupcache/lru"
)

// int:PrivateMessage
var PrivateMessageLru = lru.New(512)

// int:GroupMessage
var GroupMessageLru = lru.New(2048)

// string:
var FriendRequestLru = lru.New(128)

// string:
var GroupRequestLru = lru.New(128)

// string:
var GroupInvitedRequestLru = lru.New(16)
