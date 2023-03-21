package cache

import (
	"sync"

	"github.com/golang/groupcache/lru"
)

type LruCache struct {
	*lru.Cache
	sync.Mutex
}

func NewLruCache(maxEntries int) *LruCache {
	return &LruCache{
		Cache: lru.New(maxEntries),
	}
}

func (l *LruCache) Add(key lru.Key, value interface{}) {
	l.Lock()
	l.Cache.Add(key, value)
	l.Unlock()
}
func (l *LruCache) Get(key lru.Key) (value interface{}, ok bool) {
	l.Lock()
	value, ok = l.Cache.Get(key)
	l.Unlock()
	return
}

// int:PrivateMessage
var PrivateMessageLru = NewLruCache(512)

// int:GroupMessage
var GroupMessageLru = NewLruCache(2048)

// int:ChannelMessage
var ChannelMessageLru = NewLruCache(2048)

// string:
var FriendRequestLru = NewLruCache(128)

// string:
var GroupRequestLru = NewLruCache(128)

// string:
var GroupInvitedRequestLru = NewLruCache(16)

var GuildAdminLru = NewLruCache(2048)

var GetGuildAdminTimeLru = NewLruCache(100)
