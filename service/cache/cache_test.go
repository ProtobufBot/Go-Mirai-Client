package cache

import "testing"

func TestLruCache_Add(t *testing.T) {
	lruCache := NewLruCache(3)
	lruCache.Add(10, 10)
	v, ok := lruCache.Get(10)
	t.Logf("1: %+v %+v", v, ok)
	lruCache.Add(11, 10)
	lruCache.Add(13, 10)
	lruCache.Add(14, 10)
	v, ok = lruCache.Get(10)
	t.Logf("2: %+v %+v", v, ok)
}
